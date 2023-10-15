package utils

import (
	"testing"
)

func TestColorText(t *testing.T) {
	// Test cases
    testCases := []struct {
        name     string
        text     string
        color    StdOutColor
        expected string
    }{
        {
            name:     "Red Color",
            text:     "Test Red Color",
            color:    Red,
            expected: "\033[31mTest Red Color\033[0m",
        },
        {
            name:     "Green Color",
            text:     "Test Green Color",
            color:    Green,
            expected: "\033[32mTest Green Color\033[0m",
        },
		{
            name:     "Yellow Color",
            text:     "Test Yellow Color",
            color:    Yellow,
            expected: "\033[33mTest Yellow Color\033[0m",
        },
        {
            name:     "Blue Color",
            text:     "Test Blue Color",
            color:    Blue,
            expected: "\033[34mTest Blue Color\033[0m",
        },
		{
            name:     "Purple Color",
            text:     "Test Purple Color",
            color:    Purple,
            expected: "\033[35mTest Purple Color\033[0m",
        },
        {
            name:     "Cyan Color",
            text:     "Test Cyan Color",
            color:    Cyan,
            expected: "\033[36mTest Cyan Color\033[0m",
        },
		{
            name:     "Gray Color",
            text:     "Test Gray Color",
            color:    Gray,
            expected: "\033[37mTest Gray Color\033[0m",
        },
        {
            name:     "White Color",
            text:     "Test White Color",
            color:    White,
            expected: "\033[97mTest White Color\033[0m",
        },
    }

	for _, tc := range testCases {
        result := ColorText(tc.text, tc.color)
        if result != tc.expected {
            t.Errorf("ColorText(%s, %s) = %s, want %s", tc.text, tc.color, result, tc.expected)
        }
    }
}