package clients

import (
	"api-gateway/internal/mocks"
	"api-gateway/internal/proto/tab"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateTab_Success(t *testing.T) {
	mockClient := new(mocks.MockTabGenClient)

	expectedResp := &tab.TabResponse{Tab: "some tab"}
	mockClient.
		On("GenerateTab", mock.Anything, mock.MatchedBy(func(req *tab.TabRequest) bool {
			return req.Audio.FileName == "test.wav"
		}), mock.Anything).
		Return(expectedResp, nil)

	tabGenClient = mockClient

	resp, err := GenerateTab(context.Background(), "test.wav", []byte("data"))

	require.NoError(t, err)
	require.Equal(t, expectedResp.Tab, resp.Tab)

	mockClient.AssertExpectations(t)
}

func TestGenerateTab_Error(t *testing.T) {
	mockClient := new(mocks.MockTabGenClient)

	mockClient.
		On("GenerateTab", mock.Anything, mock.AnythingOfType("*tab.TabRequest"), mock.Anything).
		Return((*tab.TabResponse)(nil), assert.AnError)

	tabGenClient = mockClient

	resp, err := GenerateTab(context.Background(), "test.wav", []byte("data"))

	require.Error(t, err)
	require.Nil(t, resp)

	mockClient.AssertExpectations(t)
}
