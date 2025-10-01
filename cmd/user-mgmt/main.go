package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	Storage storageAPI = nil
)

type CtxKey struct {
	key string
}

var (
	CtxKeyTrace = CtxKey{key: "trace"}
)

func randalphanum() string {
	v := ""
	for len(v) < 32 {
		v = fmt.Sprintf("%v%v", v, strconv.FormatInt(rand.Int63(), 16))
	}
	v = v[:32]
	return v
}

func AddTrace(ctx context.Context, w http.ResponseWriter) context.Context {
	trace := w.Header().Get("trace")
	if len(trace) == 0 {
		trace = randalphanum()
		w.Header().Add("trace", trace)
	}
	if ctx.Value(CtxKeyTrace) == nil {
		ctx = context.WithValue(ctx, CtxKeyTrace, trace)
	}
	return ctx
}

func GetTrace(w http.ResponseWriter) string {
	if len(w.Header().Get("trace")) == 0 {
		return ""
	}
	return w.Header().Get("trace")
}

var Spans map[string][]*Span = make(map[string][]*Span)

type Span struct {
	FuncName string    `json:"func_name,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Line     int       `json:"line,omitempty"`
	Trace    string    `json:"trace,omitempty"`
	Start    time.Time `json:"start,omitempty"`
	End      time.Time `json:"end,omitempty"`
	Desc     string    `json:"desc,omitempty"`
}

func PrintSpans(ctx context.Context) {
	trace := ctx.Value(CtxKeyTrace).(string)
	spans := Spans[trace]
	emptyTime := time.Time{}
	var prevEnd time.Time
	for _, span := range spans {
		if span.End.Equal(emptyTime) {
			if prevEnd.Equal(emptyTime) {
				span.End = time.Now().UTC()
			} else {
				span.End = prevEnd
			}
			prevEnd = span.End
		}
		bs, err := json.Marshal(span)
		if err == nil {
			fmt.Println(string(bs))
		}
	}
}

func StartSpan(ctx context.Context, desc string) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	trace := ctx.Value(CtxKeyTrace).(string)
	s := &Span{
		FuncName: fn.Name(),
		Filename: file,
		Line:     line,
		Trace:    trace,
		Start:    time.Now().UTC(),
		Desc:     desc,
	}
	Spans[trace] = append(Spans[trace], s)
}

func StopSpan(ctx context.Context, desc string) {
	pc, file, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	trace := ctx.Value(CtxKeyTrace).(string)
	spans, ok := Spans[trace]
	if !ok {
		return
	}
	for _, span := range spans {
		if span.Filename == file && span.FuncName == fn.Name() && span.Desc == desc {
			span.End = time.Now().UTC()
			break
		}
	}
}

type ierror interface {
	String() string
	StatusCode() int
	StatusMessage() string
}

type erro struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *erro) String() string {
	return fmt.Sprintf(`{"code": %d, "message": %q}`, e.Code, e.Message)
}

func (e *erro) StatusCode() int {
	return e.Code
}
func (e *erro) StatusMessage() string {
	return e.Message
}

func NewError(code uint, format string, a ...any) ierror {
	return &erro{
		Code:    int(code),
		Message: fmt.Sprintf(format, a...),
	}
}

type Meta struct {
	Version *uint      `json:"version,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	Updated *time.Time `json:"updated,omitempty"`
	Deleted *bool      `json:"deleted,omitempty"`
}

type User struct {
	ID         *uint   `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Email      *string `json:"email,omitempty"`
	Supervisor *string `json:"supervisor,omitempty"`
	Meta       *Meta   `json:"meta,omitempty"`
}

type Action int

const (
	CREATE Action = iota
	UPDATE
	SOFTDELETE
	HARDDELETE
	READ
	READALL
	READVERSIONS
	READSUBORDINATES
)

type storageAPI interface {
	Close(ctx context.Context)
	Read(ctx context.Context, action Action, id uint) ([]*User, ierror)
	Write(ctx context.Context, action Action, user *User) (*User, ierror)
}

type LocalStorage map[uint][]*User

func (s LocalStorage) Close(ctx context.Context) {

	StartSpan(ctx, "LocalStorage:Close")
	defer StopSpan(ctx, "LocalStorage:Close")

	s = nil
}

func (s LocalStorage) Read(ctx context.Context, action Action, id uint) ([]*User, ierror) {
	StartSpan(ctx, "LocalStorage:Read")
	defer StopSpan(ctx, "LocalStorage:Read")

	if s == nil {
		return nil, NewError(http.StatusInternalServerError, "localstorage not initialized for read")
	}

	switch action {
	case READ:
		StartSpan(ctx, "LocalStorage:Read:READ")
		defer StopSpan(ctx, "LocalStorage:Read:READ")
		us, ok := s[id]
		if !ok || len(us) == 0 {
			return nil, NewError(http.StatusNotFound, "not found during read")
		}
		return []*User{us[len(us)-1]}, nil
	case READVERSIONS:
		StartSpan(ctx, "LocalStorage:Read:READVERSIONS")
		defer StopSpan(ctx, "LocalStorage:Read:READVERSIONS")
		us, ok := s[id]
		if !ok || len(us) == 0 {
			return nil, NewError(http.StatusNotFound, "not found during read")
		}
		return us, nil
	case READALL:
		StartSpan(ctx, "LocalStorage:Read:READALL")
		defer StopSpan(ctx, "LocalStorage:Read:READALL")
		uss := make([]*User, len(s))
		i := 0
		for _, us := range s {
			if len(us) == 0 {
				continue
			}
			uss[i] = us[len(us)-1]
			i += 1
		}
		return uss, nil
	default:
		return nil, NewError(http.StatusInternalServerError, "undefined read action")
	}
}

func (s LocalStorage) Write(ctx context.Context, action Action, user *User) (*User, ierror) {

	StartSpan(ctx, "LocalStorage:Write")
	defer StopSpan(ctx, "LocalStorage:Write")

	if s == nil {
		return nil, NewError(http.StatusInternalServerError, "localstorage not initialized for write")
	}
	if user == nil {
		return nil, NewError(http.StatusInternalServerError, "invalid user object in write")
	}

	var us []*User
	var ok bool
	if user.ID != nil {
		us, ok = s[*user.ID]
	}

	switch action {

	case CREATE:
		StartSpan(ctx, "LocalStorage:Write:CREATE")
		defer StopSpan(ctx, "LocalStorage:Write:CREATE")
		if ok || len(us) > 0 {
			return nil, NewError(http.StatusConflict, "already exists during write")
		}
		id := uint(len(s) + 1)
		user.ID = &id
		s[id] = append(s[id], user)
		return user, nil

	case UPDATE:
		StartSpan(ctx, "LocalStorage:Write:UPDATE")
		defer StopSpan(ctx, "LocalStorage:Write:UPDATE")
		if !ok || len(us) == 0 {
			return nil, NewError(http.StatusNotFound, "not found during write")
		}
		var updUser User = *(us[len(us)-1])
		updated := false
		if user.Name != nil {
			updUser.Name = user.Name
			updated = true
		}
		if user.Email != nil {
			updUser.Email = user.Email
			updated = true
		}
		if user.Supervisor != nil {
			updUser.Supervisor = user.Supervisor
			updated = true
		}
		if updated {
			now := time.Now().UTC()
			updUser.Meta.Updated = &now
			*updUser.Meta.Version += 1
			s[*user.ID] = append(s[*user.ID], &updUser)
		}
		us = s[*user.ID]
		return us[len(us)-1], nil

	case SOFTDELETE:
		StartSpan(ctx, "LocalStorage:Write:SOFTDELETE")
		defer StopSpan(ctx, "LocalStorage:Write:SOFTDELETE")
		if ok {
			for _, u := range us {
				b := true
				u.Meta.Deleted = &b
			}
		}
		return nil, nil

	case HARDDELETE:
		StartSpan(ctx, "LocalStorage:Write:HARDDELETE")
		defer StopSpan(ctx, "LocalStorage:Write:HARDDELETE")
		if ok {
			delete(s, *user.ID)
		}
		return nil, nil

	default:
		return nil, NewError(http.StatusInternalServerError, "undefined write action")
	}
}

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	ctx := AddTrace(context.Background(), w)
	defer PrintSpans(ctx)

	StartSpan(ctx, "handlePostUser")
	defer StopSpan(ctx, "handlePostUser")

	var user User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		err := NewError(http.StatusBadRequest, "invalid payload")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	{
		// validation
		if user.ID != nil {
			err := NewError(http.StatusBadRequest, "id not allowed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	u, err := Storage.Write(ctx, CREATE, &user)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := AddTrace(context.Background(), w)
	defer PrintSpans(ctx)

	StartSpan(ctx, "handleGetUserByID")
	defer StopSpan(ctx, "handleGetUserByID")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	u, err := Storage.Read(ctx, READ, uint(id))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := AddTrace(context.Background(), w)
	defer PrintSpans(ctx)

	StartSpan(ctx, "handleGetUsers")
	defer StopSpan(ctx, "handleGetUsers")

	us, err := Storage.Read(ctx, READALL, 0)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(us)
	if jsonerr != nil {
		err := NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}

// TODO: manual test
func handleGetUserVersions(w http.ResponseWriter, r *http.Request) {
	ctx := AddTrace(context.Background(), w)
	defer PrintSpans(ctx)

	StartSpan(ctx, "handleGetUserVersions")
	defer StopSpan(ctx, "handleGetUserVersions")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	u, err := Storage.Read(ctx, READVERSIONS, uint(id))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}

func handleGetUserSubOrdinates(w http.ResponseWriter, r *http.Request) {}

func handlePatchUser(w http.ResponseWriter, r *http.Request) {}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {}

func main() {

	Storage = make(LocalStorage)

	http.HandleFunc("POST /user", handlePostUser)
	http.HandleFunc("GET /user/{id}", handleGetUserByID)
	http.HandleFunc("GET /users", handleGetUsers)
	http.HandleFunc("GET /user/{id}/versions", handleGetUserVersions)
	http.HandleFunc("GET /user/{id}/subordinates", handleGetUserSubOrdinates)
	http.HandleFunc("PATCH /user", handlePatchUser)
	http.HandleFunc("DELETE /user/{id}", handleDeleteUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
