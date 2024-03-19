package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
	"vk-rest/pkg/middleware"
	"vk-rest/pkg/models"
	httpResponse "vk-rest/pkg/response"
	"vk-rest/usecase"
)

type Api struct {
	log  *logrus.Logger
	mx   *http.ServeMux
	core usecase.ICore
}

func GetApi(core *usecase.Core, log *logrus.Logger) *Api {
	api := &Api{
		core: core,
		log:  log,
		mx:   http.NewServeMux(),
	}

	api.mx.HandleFunc("/signin", api.Signin)
	api.mx.HandleFunc("/signup", api.Signup)
	api.mx.HandleFunc("/logout", api.Logout)
	api.mx.HandleFunc("/authcheck", api.AuthAccept)

	api.mx.Handle("/api/v1/questions/add", middleware.AuthCheck(http.HandlerFunc(api.QuestionAdd), core, log))
	api.mx.Handle("/api/v1/questions/event", middleware.AuthCheck(http.HandlerFunc(api.QuestionEvent), core, log))
	api.mx.Handle("/api/v1/questions/user", middleware.AuthCheck(http.HandlerFunc(api.QuestionUser), core, log))

	return api
}

func (a *Api) ListenAndServe(port string) error {
	err := http.ListenAndServe(":"+port, a.mx)
	if err != nil {
		a.log.Error("ListenAndServer error: ", err.Error())
		return err
	}

	return nil
}

func (a *Api) QuestionAdd(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.Quest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Read body in question add error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Unmarshal in question add error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, err = a.core.QuestionAdd(&request)
	if err != nil {
		a.log.Error("Question add error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) QuestionEvent(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.EventItem

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Read body in question event error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Unmarshal in question event error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.QuestionEvent(&request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) QuestionUser(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary signIn
// @Tags Auth
// @Description authenticate user by providing login and password credentials
// @ID authenticate-user
// @Accept json
// @Produce json
// @Param input body models.SigninRequest false "login and password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /signin [post]
func (a *Api) Signin(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var authorized bool

	sessionCookie, err := r.Cookie("session_id")
	if err == nil && sessionCookie != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), sessionCookie.Value)
	}

	if authorized {
		response.Status = http.StatusConflict
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, found, err := a.core.FindUserAccount(request.Login, request.Password)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if !found {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, _ := a.core.CreateSession(r.Context(), request.Login)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.SID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary signUp
// @Tags Auth
// @Desription create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body models.SignupRequest false "account information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /signup [post]
func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.SignupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	found, err := a.core.FindUserByLogin(request.Login)
	if err != nil {
		a.log.Error("Signup error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if found {
		response.Status = http.StatusConflict
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.CreateUserAccount(request.Login, request.Password)
	if err != nil {
		a.log.Error("create user error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary end current user session
// @Tags Auth
// @ID logout
// @Produce json
// @Param session_id header string false "Session ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /logout [delete]
func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodDelete {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.KillSession(r.Context(), cookie.Value)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @summary check authentication status and return user info
// @description returns user info if they are currently logged in
// @Tags Auth
// @produce application/json
// @Param session_id header string false "Session ID"
// @success 200 {object} models.AuthCheckResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @router /authcheck [get]
func (a *Api) AuthAccept(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var authorized bool

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), session.Value)
	}
	a.log.Warning("API", authorized)
	if !authorized {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	login, err := a.core.GetUserName(r.Context(), session.Value)
	if err != nil {
		a.log.Error("auth accept error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = models.AuthCheckResponse{
		Login: login,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}
