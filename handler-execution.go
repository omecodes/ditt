package ditt

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"sync"
)

type handlerExecution struct {
	BaseHandler
}

func (e *handlerExecution) Login(ctx context.Context, login string, password string) (bool, error) {
	if login == "admin" {
		return password == Env.AdminPassword, nil
	}

	userData, err := e.GetUser(ctx, login)
	if err != nil {
		if err == NotFound {
			return false, NotAuthorized
		}
		return false, err
	}

	hashedPassword := userData.Password()
	return err == bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)), nil
}

func (e *handlerExecution) AddUsers(_ context.Context, reader io.Reader) error {

	tasksResultsChannelSignal := make(chan chan UserDataProcessingResult)
	defer close(tasksResultsChannelSignal)

	runResultChannelSignal := make(chan chan error)
	defer close(runResultChannelSignal)

	processor := func(data UserData) (UserData, error) {
		processedData, err := processData(writeProcessors, data)
		if err != nil {
			return "", err
		}
		return "", Env.DataStore.Save(processedData)
	}

	runner := ConcurrentUserDataProcessingRunner{
		Provider:            newJsonObjectStreamParser(reader).parseUsers,
		Processor:           UserDataProcessorFunc(processor),
		TasksResultsSignals: tasksResultsChannelSignal,
		ResultSignal:        runResultChannelSignal,
	}
	runner.Run()

	tasksResults := <-tasksResultsChannelSignal
	for {
		result, hasMore := <-tasksResults
		if !hasMore {
			break
		}

		if result.Err != nil {
			log.Println("data", result.UserId, ":", result.Err)
		} else {
			log.Println("data", result.UserId, ": saved")
		}
	}

	runResult := <-runResultChannelSignal
	return <-runResult
}

func (e *handlerExecution) DeleteUser(_ context.Context, userId string) error {
	err := Env.Files.Delete(userId)
	if err != nil {
		return err
	}
	return Env.DataStore.Delete(userId)
}

func (e *handlerExecution) GetUser(_ context.Context, userId string) (UserData, error) {
	userData, err := Env.DataStore.Get(userId)
	if err != nil {
		return "", err
	}

	return processData(readProcessors, userData)
}

func (e *handlerExecution) GetUserList(ctx context.Context, opts ListOptions) (*UserDataList, error) {
	tasksResultsChannelSignal := make(chan chan UserDataProcessingResult)
	defer close(tasksResultsChannelSignal)

	runResultChannelSignal := make(chan chan error)
	defer close(runResultChannelSignal)

	provider := func(callback UserDataCallback) error {
		userId := GetLoggedUser(ctx)
		var err error
		if userId == "admin" {
			err = Env.DataStore.List(opts.Offset, opts.Count, callback)
		} else {
			err = Env.DataStore.ListForUser(userId, opts.Offset, opts.Count, callback)
		}
		return err
	}
	processor := func(data UserData) (UserData, error) {
		return processData(readProcessors, data)
	}
	runner := ConcurrentUserDataProcessingRunner{
		Provider:            provider,
		Processor:           UserDataProcessorFunc(processor),
		TasksResultsSignals: tasksResultsChannelSignal,
		ResultSignal:        runResultChannelSignal,
	}
	runner.Run()

	userDataList := &UserDataList{
		Offset: opts.Offset,
	}
	tasksResults := <-tasksResultsChannelSignal
	for {
		result, hasMore := <-tasksResults
		if !hasMore {
			break
		}

		if result.Err != nil {
			log.Println("data", result.UserId, ":", result.Err)
		} else {
			//userDataList.UserDataList = append(userDataList.UserDataList, result.Data)
		}
	}

	runResult := <-runResultChannelSignal
	err := <-runResult

	log.Println("done listing:", err)
	return userDataList, err
}

func (e *handlerExecution) loadFileContent(data UserData, wg *sync.WaitGroup, processed chan<- UserData, failures chan<- error) {
	defer wg.Done()
	processedData, pErr := processData(readProcessors, data)
	if pErr != nil {
		failures <- pErr
	} else {
		processed <- processedData
	}
}

func (e *handlerExecution) UpdateUser(_ context.Context, _ string, userData UserData) error {
	processedData, err := processData(writeProcessors, userData)
	if err != nil {
		return err
	}

	err = Env.DataStore.Save(processedData)
	if err != nil {
		return err
	}

	return nil
}
