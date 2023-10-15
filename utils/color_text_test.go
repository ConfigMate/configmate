package utils

import (
	"utils"
	"testing"

	"github.com/ConfigMate/configmate/parsers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestColorText(t *testing.T) {
	// Test cases
    testCases := []struct {
        name     string
        text     string
        color    utils.StdOutColor
        expected string
    }{
        {
            name:     "Red color",
            text:     "Test Red Color",
            color:    utils.StdOutColor.Red,
            expected: "\033[31mTest Red Color\033[0m",
        },
        {
            name:     "Green color",
            text:     "Test Green Color",
            color:    utils.StdOutColor.Green,
            expected: "\033[32mTest Green Color\033[0m",
        },
		{
            name:     "Yellow color",
            text:     "Test Yellow Color",
            color:    utils.StdOutColor.Red,
            expected: "\033[33mTest Yellow Color\033[0m",
        },
        {
            name:     "Blue color",
            text:     "Test Blue Color",
            color:    utils.StdOutColor.Green,
            expected: "\033[34mTest Blue Color\033[0m",
        },
		{
            name:     "Purple color",
            text:     "Test Purple Color",
            color:    utils.StdOutColor.Red,
            expected: "\033[35mTest Purple Color\033[0m",
        },
        {
            name:     "Cyan color",
            text:     "Test Cyan Color",
            color:    utils.StdOutColor.Green,
            expected: "\033[36mTest Cyan Color\033[0m",
        },
		{
            name:     "Gray color",
            text:     "Test Gray Color",
            color:    utils.StdOutColor.Red,
            expected: "\033[37mTest Gray Color\033[0m",
        },
        {
            name:     "White color",
            text:     "Test White Color",
            color:    utils.StdOutColor.Green,
            expected: "\033[97mTest White Color\033[0m",
        },
    }

	for _, tc := range testCases {
        result := utils.ColorText(tc.text, tc.color)
        if result != tc.expected {
            t.Errorf("ColorText(%s, %s) = %s, want %s", tc.text, tc.color, result, tc.expected)
        }
    }
}