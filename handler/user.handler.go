package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"radar/common/response"
	"radar/model"
	"radar/service"
	"radar/utils"
	"strconv"
)

type UserHandler interface {
	SendVerificationMail() http.HandlerFunc
	VerifyAccount() http.HandlerFunc
	CreateEvent() http.HandlerFunc
	FilterEventsBy() http.HandlerFunc
	AllEvents()  http.HandlerFunc
	AskQuestion() http.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// SendVerificationEmail sends the verification email

func (c *userHandler) SendVerificationMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")

		_, err := c.userService.FindUser(email)
		fmt.Println("email: ", email)
		fmt.Println("err: ", err)

		if err == nil {
			err = c.userService.SendVerificationEmail(email)
		}

		fmt.Println(err)

		if err != nil {
			response := response.ErrorResponse("Error while sending verification mail", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "Verification mail sent successfully", email)
		utils.ResponseJSON(w, response)
	}
}

// verifyAccount verifies the account

func (c *userHandler) VerifyAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")
		code, _ := strconv.Atoi(r.URL.Query().Get("Code"))

		err := c.userService.VerifyAccount(email, code)

		if err != nil {
			response := response.ErrorResponse("Verification failed, Invalid OTP", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "Account verified successfully", email)
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newEvent model.Event
		json.NewDecoder(r.Body).Decode(&newEvent)
		newEvent.Organizer_name = (r.Header.Get("Organizer_name"))
		_, err := c.userService.CreateEvent(newEvent)
		if err != nil {
			response := response.ErrorResponse("Failed to add new post", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "SUCCESS", newEvent)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) FilterEventsBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		free := r.URL.Query().Get("Free")
		sex := r.URL.Query().Get("Sex")
		fmt.Println("free from handlers:",free)
		cusat_only := (r.URL.Query().Get("Cusat_only"))
		fmt.Println("cusat only from handler:",cusat_only)

		events, err := c.userService.FilterEventsBy( sex,cusat_only, free)

		result := struct {
			Events *[]model.EventResponse
		}{
			Events: events,
		}

		if err != nil {
			response := response.ErrorResponse("error while getting posts from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "All Events", result)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}


func (c *userHandler) AllEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		events, err := c.userService.AllEvents()

		result := struct {
			Events *[]model.EventResponse
		}{
			Events: events,
		}

		if err != nil {
			response := response.ErrorResponse("error while getting posts from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "All Events", result)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}


func (c *userHandler) AskQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newQuestion model.FAQA
		json.NewDecoder(r.Body).Decode(&newQuestion)

		if newQuestion.Question == "" {

				log.Fatal("Qustion box is empty!")
				return
		}
		newQuestion.User_name = r.Header.Get("User_name")
		newQuestion.Event_name = r.Header.Get("Event_name")
		err := c.userService.AskQuestion(newQuestion)
		if err != nil {
			response := response.ErrorResponse("Failed to add new comment", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "SUCCESS", newQuestion)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}
