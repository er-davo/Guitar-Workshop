package service

import (
	"api-gateway/internal/mocks"
	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTabServiceGenerateTab(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	audiosep := mocks.NewMockAudioSeparator(ctrl)
	tabgen := mocks.NewMockTabGenerator(ctrl)

	service := NewTabService(nil, nil, tabgen, audiosep)

	type testCase struct {
		name          string
		separation    bool
		setupMocks    func()
		expectedBody  string
		expectedError string
	}

	tests := []testCase{
		{
			name:       "success with separation",
			separation: true,
			setupMocks: func() {
				audiosep.EXPECT().
					SeparateAudio(gomock.Any(), "file.wav", []byte("data")).
					Return(map[string]*separator.AudioFileData{
						"other": {FileName: "other.wav", AudioBytes: []byte("guitar-data")},
					}, nil)

				tabgen.EXPECT().
					GenerateTab(gomock.Any(), "other.wav", []byte("guitar-data")).
					Return(&tab.TabResponse{Tab: "TAB123"}, nil)
			},
			expectedBody: "TAB123",
		},
		{
			name:       "success without separation",
			separation: false,
			setupMocks: func() {
				tabgen.EXPECT().
					GenerateTab(gomock.Any(), "file.wav", []byte("data")).
					Return(&tab.TabResponse{Tab: "TAB_NOSEP"}, nil)
			},
			expectedBody: "TAB_NOSEP",
		},
		{
			name:       "separation error",
			separation: true,
			setupMocks: func() {
				audiosep.EXPECT().
					SeparateAudio(gomock.Any(), "file.wav", []byte("data")).
					Return(nil, errors.New("sep fail"))
			},
			expectedError: "sep fail",
		},
		{
			name:       "missing 'other' stem",
			separation: true,
			setupMocks: func() {
				audiosep.EXPECT().
					SeparateAudio(gomock.Any(), "file.wav", []byte("data")).
					Return(map[string]*separator.AudioFileData{}, nil)
			},
			expectedError: "audio separation result missing 'other' stem",
		},
		{
			name:       "tab generator error",
			separation: true,
			setupMocks: func() {
				audiosep.EXPECT().
					SeparateAudio(gomock.Any(), "file.wav", []byte("data")).
					Return(map[string]*separator.AudioFileData{
						"other": {FileName: "other.wav", AudioBytes: []byte("guitar-data")},
					}, nil)

				tabgen.EXPECT().
					GenerateTab(gomock.Any(), "other.wav", []byte("guitar-data")).
					Return(nil, errors.New("tabgen fail"))
			},
			expectedError: "tabgen fail",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			tab, err := service.GenerateTab(context.Background(), "file.wav", []byte("data"), tc.separation)

			if tc.expectedError != "" {
				require.Nil(t, tab)
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tab)
				require.Equal(t, tc.expectedBody, tab.Body)
			}
		})
	}
}
