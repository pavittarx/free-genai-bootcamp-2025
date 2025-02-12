package models_test

import (
	"testing"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

func TestWord_Validate(t *testing.T) {
	tests := []struct {
		name    string
		word    models.Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: models.Word{
				Hindi:    "नमस्ते",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantErr: false,
		},
		{
			name: "empty hindi",
			word: models.Word{
				Hindi:    "",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantErr: true,
		},
		{
			name: "empty english",
			word: models.Word{
				Hindi:    "नमस्ते",
				English:  "",
				Hinglish: "Namaste",
			},
			wantErr: true,
		},
		{
			name: "invalid hindi characters",
			word: models.Word{
				Hindi:    "नमस्ते123",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantErr: true,
		},
		{
			name: "invalid english characters",
			word: models.Word{
				Hindi:    "नमस्ते",
				English:  "Hello123!",
				Hinglish: "Namaste",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.word.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Word.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWord_Sanitize(t *testing.T) {
	word := &models.Word{
		Hindi:    "  नमस्ते  ",
		English:  "  Hello  ",
		Hinglish: "  Namaste  ",
	}

	word.Sanitize()

	if word.Hindi != "नमस्ते" {
		t.Errorf("Word.Sanitize() Hindi = %v, want %v", word.Hindi, "नमस्ते")
	}
	if word.English != "Hello" {
		t.Errorf("Word.Sanitize() English = %v, want %v", word.English, "Hello")
	}
	if word.Hinglish != "Namaste" {
		t.Errorf("Word.Sanitize() Hinglish = %v, want %v", word.Hinglish, "Namaste")
	}
}

func TestWord_GenerateScrambledWord(t *testing.T) {
	word := &models.Word{
		Hindi: "नमस्ते",
	}

	word.GenerateScrambledWord()

	if word.Scrambled == "" {
		t.Error("Word.GenerateScrambledWord() failed to generate scrambled word")
	}
	if word.Scrambled == word.Hindi {
		t.Error("Word.GenerateScrambledWord() generated same word as original")
	}
	if len(word.Scrambled) != len(word.Hindi) {
		t.Errorf("Word.GenerateScrambledWord() length mismatch: got %v, want %v", len(word.Scrambled), len(word.Hindi))
	}
}
