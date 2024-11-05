package repository

// Мы предполагаем, что у вас настроен mock-объект для Redis
// type MockRedisConn struct {
// 	mock *gomock.Controller
// }

// func (m *MockRedisConn) Conn(ctx context.Context) redis.Cmdable {
// 	mockRedisClient := NewMockCmdable(m.mock)
// 	return mockRedisClient
// }

// // Mock Cmdable methods implementation generated with Mockgen
// type MockCmdable struct {
// 	mock.Mock
// }

// func (m *MockCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
// 	args := m.Called(ctx, key, value, expiration)
// 	return args.Get(0).(*redis.StatusCmd)
// }

// // Mock generation is typically done with mockgen tool

// func TestRegisterSessionRedis(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Create a mock redis client
// 	mockRedis := NewMockCmdable(ctrl)
// 	mockRedis.EXPECT().Set(gomock.Any(), "test_cookie", 1, 7*24*time.Hour).Return(redis.NewStatusResult("", nil))

// 	authRepo := &AuthRepository{
// 		redisDb: &MockRedisConn{mock: ctrl},
// 	}

// 	err := authRepo.RegisterSessionRedis("test_cookie", 1)
// 	assert.Nil(t, err, fmt.Sprintf("Expected no error, but got %v", err))
// }

// Здесь мы определяем необходимые поля в структуре UserProfile

// // Ошибка для неверных учетных данных
// var ErrWrongCredentials = errors.New("wrong credentials")

// func TestGetUserByEmail(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

// 	repo := &AuthRepository{
// 		db: mockPool,
// 	}

// 	// Определяем переменные для теста
// 	email := "test@example.com"
// 	expectedUser := &models.UserProfile{
// 		ID:           1,
// 		Name:         "Test User",
// 		Email:        email,
// 		Description:  "Test Description",
// 		JoinedAt:     time.Now(),
// 		UpdatedAt:    time.Now(),
// 		PasswordHash: "testhash",
// 	}

// 	// Создаем ожидаемый response от базы данных
// 	rows := pgxpoolmock.NewRows([]string{"u_id", "nickname", "email", "description", "joined_at", "updated_at", "password_hash"}).
// 		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Description, expectedUser.JoinedAt, expectedUser.UpdatedAt, expectedUser.PasswordHash)

// 	mockPool.EXPECT().
// 		QueryRow(gomock.Any(), gomock.Any(), email).
// 		Return(rows)

// 	// Выполнение тестируемой функции
// 	user, err := repo.GetUserByEmail(email)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedUser, user)

// 	// Проверка на случай, если пользователь не найден
// 	mockPool.EXPECT().
// 		QueryRow(gomock.Any(), gomock.Any(), "nonexistent@example.com").
// 		Return(pgxpoolmock.NewRows(nil)) // No rows

// 	user, err = repo.GetUserByEmail("nonexistent@example.com")
// 	assert.NotNil(t, err)
// 	assert.Equal(t, ErrWrongCredentials, err)
// 	assert.Nil(t, user)
// }
