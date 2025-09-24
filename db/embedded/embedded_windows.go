//go:build windows

package embedded

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"r3/config"
	"r3/log"
	"r3/tools"
	"strings"
	"syscall"
	"time"
)

func Start() error {

	// check for existing embedded database path
	exists, err := tools.Exists(dbData)
	if err != nil {
		return err
	}
	if !exists {
		templatePath := strings.Replace(dbData, "database", "database_template", 1)
		templateExists, err := tools.Exists(templatePath)
		if err != nil {
			return err
		}

		if templateExists {
			// get database from template
			if err := tools.FileMove(templatePath, dbData, false); err != nil {
				return err
			}
		} else {
			// initialize database using initdb if template doesn't exist
			if err := initializeDatabase(); err != nil {
				return err
			}
		}
	}

	// check embedded database state
	state, err := status()
	if err != nil {
		return err
	}

	if state {
		return fmt.Errorf("database already running, another instance is likely active")
	}
	_, err = execWaitFor(dbBinCtl, []string{"start", "-D", dbData,
		fmt.Sprintf(`-o "-p %d"`, config.File.Db.Port)}, []string{msgStarted}, 10)

	return err
}

func Stop() error {

	state, err := status()
	if err != nil {
		return err
	}

	if !state {
		log.Info(log.ContextServer, "embedded database already stopped")
		return nil
	}

	_, err = execWaitFor(dbBinCtl, []string{"stop", "-D", dbData}, []string{msgStopped}, 10)
	return err
}

func status() (bool, error) {

	foundLine, err := execWaitFor(dbBinCtl, []string{"status", "-D", dbData},
		[]string{msgState0, msgState1}, 5)

	if err != nil {
		return false, err
	}
	// returns true if DB server is running
	return strings.Contains(foundLine, msgState1), nil
}

// initializeDatabase creates a new PostgreSQL database using initdb
func initializeDatabase() error {
	// Create initdb command
	cmd := exec.Command(filepath.Join(dbBin, "initdb"),
		"-D", dbData,
		"-U", config.File.Db.User,
		"--auth-local=trust", 
		"--auth-host=md5",
		"-E", "UTF8",
		"--locale=C",
		"--no-instructions")
	
	tools.CmdAddSysProgAttrs(cmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("LC_MESSAGES=%s", locale))
	
	// Set up separate process group for clean shutdown
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	
	// Run initdb and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v, output: %s", err, string(output))
	}
	
	log.Info(log.ContextServer, "Database cluster initialized successfully with initdb")
	
	// Now we need to start the database temporarily to create the application database and user
	if err := createApplicationDatabase(); err != nil {
		return fmt.Errorf("failed to create application database and user: %v", err)
	}
	if err := createApplicationUser(); err != nil {
		return fmt.Errorf("failed to create application database and user: %v", err)
	}
	
	return nil
}

// createApplicationDatabaseAndUser creates the application database and ensures proper user setup
func createApplicationDatabase() error {
	// Start the database server temporarily to create the application database
	_, err := execWaitFor(dbBinCtl, []string{"start", "-D", dbData,
		fmt.Sprintf(`-o "-p %d"`, config.File.Db.Port)}, []string{msgStarted}, 30)
	if err != nil {
		return fmt.Errorf("failed to start database for setup: %v", err)
	}
	
	// Ensure we stop the database when done, even if there's an error
	defer func() {
		execWaitFor(dbBinCtl, []string{"stop", "-D", dbData}, []string{msgStopped}, 10)
	}()
	
	// Create the application database and user using psql
	createDbCmd := exec.Command(filepath.Join(dbBin, "psql"),
		"-p", fmt.Sprintf("%d", config.File.Db.Port),
		"-U", config.File.Db.User,
		"-d", "postgres", // Connect to the default postgres database first
		"-c", "CREATE DATABASE app WITH OWNER = app TEMPLATE = template0 ENCODING = 'UTF8';"
	
	tools.CmdAddSysProgAttrs(createDbCmd)
	createDbCmd.Env = append(os.Environ(), fmt.Sprintf("LC_MESSAGES=%s", locale))
	
	output, err := createDbCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create application database: %v, output: %s", err, string(output))
	}
	
	log.Info(log.ContextServer, fmt.Sprintf("Application database '%s' created successfully", config.File.Db.Name))
	return nil
}
// createApplicationDatabaseAndUser creates the application database and ensures proper user setup
func createApplicationUser() error {
	// Start the database server temporarily to create the application database
	_, err := execWaitFor(dbBinCtl, []string{"start", "-D", dbData,
		fmt.Sprintf(`-o "-p %d"`, config.File.Db.Port)}, []string{msgStarted}, 30)
	if err != nil {
		return fmt.Errorf("failed to start database for setup: %v", err)
	}
	
	// Ensure we stop the database when done, even if there's an error
	defer func() {
		execWaitFor(dbBinCtl, []string{"stop", "-D", dbData}, []string{msgStopped}, 10)
	}()
	
	// Create the application database and user using psql
	createDbCmd := exec.Command(filepath.Join(dbBin, "psql"),
		"-p", fmt.Sprintf("%d", config.File.Db.Port),
		"-U", config.File.Db.User,
		"-d", "postgres", // Connect to the default postgres database first
		"-c", "CREATE ROLE app WITH LOGIN PASSWORD 'app!';"
	
	tools.CmdAddSysProgAttrs(createDbCmd)
	createDbCmd.Env = append(os.Environ(), fmt.Sprintf("LC_MESSAGES=%s", locale))
	
	output, err := createDbCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create application database: %v, output: %s", err, string(output))
	}
	
	log.Info(log.ContextServer, fmt.Sprintf("Application database '%s' created successfully", config.File.Db.Name))
	return nil
}

// executes call and waits for specified lines to return
// will return automatically after timeout
func execWaitFor(call string, args []string, waitFor []string, waitTime int) (string, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(waitTime)*time.Second)
	cmd := exec.CommandContext(ctx, call, args...)
	tools.CmdAddSysProgAttrs(cmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("LC_MESSAGES=%s", locale))

	// create as separate process for clean shutdown, otherwise child progs are killed immediately on SIGINT
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	type chanReturnType struct {
		err  error
		line string
	}
	chanReturn := make(chan chanReturnType)

	// react to call timeout
	go func() {
		for {
			<-ctx.Done()
			chanReturn <- chanReturnType{err: errors.New("timeout reached")}
			return
		}
	}()

	// react to new lines from standard output
	go func() {
		if err := cmd.Start(); err != nil {
			chanReturn <- chanReturnType{err: err}
			return
		}

		buf := bufio.NewReader(stdout)
		bufLines := []string{}
		for {
			line, _, err := buf.ReadLine()
			if err != nil {
				if err != io.EOF {
					// log error if not end-of-file
					log.Error(log.ContextServer, "failed to read from std out", err)
				}
				break
			}
			bufLines = append(bufLines, string(line))

			// success if expected lines turned up
			for _, waitLine := range waitFor {
				if strings.Contains(string(line), waitLine) {
					chanReturn <- chanReturnType{
						err:  nil,
						line: waitLine,
					}
					return
				}
			}
		}

		if len(bufLines) == 0 {
			// nothing turned up
			chanReturn <- chanReturnType{err: errors.New("output is empty")}
		} else {
			// expected lines did not turn up
			chanReturn <- chanReturnType{err: fmt.Errorf("unexpected output, got: %s", strings.Join(bufLines, ","))}
		}
	}()

	res := <-chanReturn
	return res.line, res.err
}
