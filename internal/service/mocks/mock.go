// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/samuraivf/twitter-clone/internal/dto"
	models "github.com/samuraivf/twitter-clone/internal/repository/models"
	service "github.com/samuraivf/twitter-clone/internal/service"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *MockAuthorization) GenerateAccessToken(username string, userId uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", username, userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockAuthorizationMockRecorder) GenerateAccessToken(username, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockAuthorization)(nil).GenerateAccessToken), username, userId)
}

// GenerateRefreshToken mocks base method.
func (m *MockAuthorization) GenerateRefreshToken(username string, userId uint) (*service.RefreshTokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", username, userId)
	ret0, _ := ret[0].(*service.RefreshTokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockAuthorizationMockRecorder) GenerateRefreshToken(username, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockAuthorization)(nil).GenerateRefreshToken), username, userId)
}

// ParseAccessToken mocks base method.
func (m *MockAuthorization) ParseAccessToken(accessToken string) (*service.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAccessToken", accessToken)
	ret0, _ := ret[0].(*service.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseAccessToken indicates an expected call of ParseAccessToken.
func (mr *MockAuthorizationMockRecorder) ParseAccessToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAccessToken", reflect.TypeOf((*MockAuthorization)(nil).ParseAccessToken), accessToken)
}

// ParseRefreshToken mocks base method.
func (m *MockAuthorization) ParseRefreshToken(refreshToken string) (*service.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseRefreshToken", refreshToken)
	ret0, _ := ret[0].(*service.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseRefreshToken indicates an expected call of ParseRefreshToken.
func (mr *MockAuthorizationMockRecorder) ParseRefreshToken(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseRefreshToken", reflect.TypeOf((*MockAuthorization)(nil).ParseRefreshToken), refreshToken)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// AddImage mocks base method.
func (m *MockUser) AddImage(image string, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImage", image, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImage indicates an expected call of AddImage.
func (mr *MockUserMockRecorder) AddImage(image, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImage", reflect.TypeOf((*MockUser)(nil).AddImage), image, userId)
}

// CreateUser mocks base method.
func (m *MockUser) CreateUser(user dto.CreateUserDto) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), user)
}

// EditProfile mocks base method.
func (m *MockUser) EditProfile(user dto.EditUserDto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditProfile indicates an expected call of EditProfile.
func (mr *MockUserMockRecorder) EditProfile(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockUser)(nil).EditProfile), user)
}

// GetUserByEmail mocks base method.
func (m *MockUser) GetUserByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUser)(nil).GetUserByEmail), email)
}

// GetUserByUsername mocks base method.
func (m *MockUser) GetUserByUsername(username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserMockRecorder) GetUserByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUser)(nil).GetUserByUsername), username)
}

// GetUserMessages mocks base method.
func (m *MockUser) GetUserMessages(userId uint) ([]*models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserMessages", userId)
	ret0, _ := ret[0].([]*models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserMessages indicates an expected call of GetUserMessages.
func (mr *MockUserMockRecorder) GetUserMessages(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserMessages", reflect.TypeOf((*MockUser)(nil).GetUserMessages), userId)
}

// Subscribe mocks base method.
func (m *MockUser) Subscribe(subscriberId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", subscriberId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockUserMockRecorder) Subscribe(subscriberId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockUser)(nil).Subscribe), subscriberId, userId)
}

// Unsubscribe mocks base method.
func (m *MockUser) Unsubscribe(subscriberId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", subscriberId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockUserMockRecorder) Unsubscribe(subscriberId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockUser)(nil).Unsubscribe), subscriberId, userId)
}

// ValidateUser mocks base method.
func (m *MockUser) ValidateUser(username, password string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUser", username, password)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockUserMockRecorder) ValidateUser(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockUser)(nil).ValidateUser), username, password)
}

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// DeleteRefreshToken mocks base method.
func (m *MockRedis) DeleteRefreshToken(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockRedisMockRecorder) DeleteRefreshToken(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockRedis)(nil).DeleteRefreshToken), ctx, key)
}

// GetRefreshToken mocks base method.
func (m *MockRedis) GetRefreshToken(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshToken", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshToken indicates an expected call of GetRefreshToken.
func (mr *MockRedisMockRecorder) GetRefreshToken(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshToken", reflect.TypeOf((*MockRedis)(nil).GetRefreshToken), ctx, key)
}

// SetRefreshToken mocks base method.
func (m *MockRedis) SetRefreshToken(ctx context.Context, key, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRefreshToken", ctx, key, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRefreshToken indicates an expected call of SetRefreshToken.
func (mr *MockRedisMockRecorder) SetRefreshToken(ctx, key, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRefreshToken", reflect.TypeOf((*MockRedis)(nil).SetRefreshToken), ctx, key, refreshToken)
}

// MockTweet is a mock of Tweet interface.
type MockTweet struct {
	ctrl     *gomock.Controller
	recorder *MockTweetMockRecorder
}

// MockTweetMockRecorder is the mock recorder for MockTweet.
type MockTweetMockRecorder struct {
	mock *MockTweet
}

// NewMockTweet creates a new mock instance.
func NewMockTweet(ctrl *gomock.Controller) *MockTweet {
	mock := &MockTweet{ctrl: ctrl}
	mock.recorder = &MockTweetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTweet) EXPECT() *MockTweetMockRecorder {
	return m.recorder
}

// CreateTweet mocks base method.
func (m *MockTweet) CreateTweet(tweetDto dto.CreateTweetDto) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTweet", tweetDto)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTweet indicates an expected call of CreateTweet.
func (mr *MockTweetMockRecorder) CreateTweet(tweetDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTweet", reflect.TypeOf((*MockTweet)(nil).CreateTweet), tweetDto)
}

// DeleteTweet mocks base method.
func (m *MockTweet) DeleteTweet(tweetId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTweet", tweetId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTweet indicates an expected call of DeleteTweet.
func (mr *MockTweetMockRecorder) DeleteTweet(tweetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTweet", reflect.TypeOf((*MockTweet)(nil).DeleteTweet), tweetId)
}

// GetTweetById mocks base method.
func (m *MockTweet) GetTweetById(id uint) (*models.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTweetById", id)
	ret0, _ := ret[0].(*models.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTweetById indicates an expected call of GetTweetById.
func (mr *MockTweetMockRecorder) GetTweetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTweetById", reflect.TypeOf((*MockTweet)(nil).GetTweetById), id)
}

// GetUserTweets mocks base method.
func (m *MockTweet) GetUserTweets(userId uint) ([]models.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTweets", userId)
	ret0, _ := ret[0].([]models.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTweets indicates an expected call of GetUserTweets.
func (mr *MockTweetMockRecorder) GetUserTweets(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTweets", reflect.TypeOf((*MockTweet)(nil).GetUserTweets), userId)
}

// LikeTweet mocks base method.
func (m *MockTweet) LikeTweet(tweetId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeTweet", tweetId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeTweet indicates an expected call of LikeTweet.
func (mr *MockTweetMockRecorder) LikeTweet(tweetId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeTweet", reflect.TypeOf((*MockTweet)(nil).LikeTweet), tweetId, userId)
}

// UnlikeTweet mocks base method.
func (m *MockTweet) UnlikeTweet(tweetId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeTweet", tweetId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeTweet indicates an expected call of UnlikeTweet.
func (mr *MockTweetMockRecorder) UnlikeTweet(tweetId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeTweet", reflect.TypeOf((*MockTweet)(nil).UnlikeTweet), tweetId, userId)
}

// UpdateTweet mocks base method.
func (m *MockTweet) UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTweet", tweetDto)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTweet indicates an expected call of UpdateTweet.
func (mr *MockTweetMockRecorder) UpdateTweet(tweetDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTweet", reflect.TypeOf((*MockTweet)(nil).UpdateTweet), tweetDto)
}

// MockComment is a mock of Comment interface.
type MockComment struct {
	ctrl     *gomock.Controller
	recorder *MockCommentMockRecorder
}

// MockCommentMockRecorder is the mock recorder for MockComment.
type MockCommentMockRecorder struct {
	mock *MockComment
}

// NewMockComment creates a new mock instance.
func NewMockComment(ctrl *gomock.Controller) *MockComment {
	mock := &MockComment{ctrl: ctrl}
	mock.recorder = &MockCommentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComment) EXPECT() *MockCommentMockRecorder {
	return m.recorder
}

// CreateComment mocks base method.
func (m *MockComment) CreateComment(commentDto dto.CreateCommentDto) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", commentDto)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockCommentMockRecorder) CreateComment(commentDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockComment)(nil).CreateComment), commentDto)
}

// DeleteComment mocks base method.
func (m *MockComment) DeleteComment(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockCommentMockRecorder) DeleteComment(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockComment)(nil).DeleteComment), id)
}

// GetCommentById mocks base method.
func (m *MockComment) GetCommentById(id uint) (*models.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentById", id)
	ret0, _ := ret[0].(*models.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentById indicates an expected call of GetCommentById.
func (mr *MockCommentMockRecorder) GetCommentById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentById", reflect.TypeOf((*MockComment)(nil).GetCommentById), id)
}

// LikeComment mocks base method.
func (m *MockComment) LikeComment(commentId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeComment", commentId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeComment indicates an expected call of LikeComment.
func (mr *MockCommentMockRecorder) LikeComment(commentId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeComment", reflect.TypeOf((*MockComment)(nil).LikeComment), commentId, userId)
}

// UnlikeComment mocks base method.
func (m *MockComment) UnlikeComment(commentId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeComment", commentId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeComment indicates an expected call of UnlikeComment.
func (mr *MockCommentMockRecorder) UnlikeComment(commentId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeComment", reflect.TypeOf((*MockComment)(nil).UnlikeComment), commentId, userId)
}

// UpdateComment mocks base method.
func (m *MockComment) UpdateComment(commentDto dto.UpdateCommentDto) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", commentDto)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComment indicates an expected call of UpdateComment.
func (mr *MockCommentMockRecorder) UpdateComment(commentDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockComment)(nil).UpdateComment), commentDto)
}

// MockTag is a mock of Tag interface.
type MockTag struct {
	ctrl     *gomock.Controller
	recorder *MockTagMockRecorder
}

// MockTagMockRecorder is the mock recorder for MockTag.
type MockTagMockRecorder struct {
	mock *MockTag
}

// NewMockTag creates a new mock instance.
func NewMockTag(ctrl *gomock.Controller) *MockTag {
	mock := &MockTag{ctrl: ctrl}
	mock.recorder = &MockTagMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTag) EXPECT() *MockTagMockRecorder {
	return m.recorder
}

// GetTagByIdWithTweets mocks base method.
func (m *MockTag) GetTagByIdWithTweets(id uint) (*models.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByIdWithTweets", id)
	ret0, _ := ret[0].(*models.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByIdWithTweets indicates an expected call of GetTagByIdWithTweets.
func (mr *MockTagMockRecorder) GetTagByIdWithTweets(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByIdWithTweets", reflect.TypeOf((*MockTag)(nil).GetTagByIdWithTweets), id)
}

// GetTagByName mocks base method.
func (m *MockTag) GetTagByName(name string) (*models.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByName", name)
	ret0, _ := ret[0].(*models.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByName indicates an expected call of GetTagByName.
func (mr *MockTagMockRecorder) GetTagByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByName", reflect.TypeOf((*MockTag)(nil).GetTagByName), name)
}
