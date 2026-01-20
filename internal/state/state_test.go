package state

import "testing"

func TestStateValid(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{Thinking, true},
		{Shaping, true},
		{Building, true},
		{Shipped, true},
		{State("invalid"), false},
		{State(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.Valid(); got != tt.want {
				t.Errorf("State(%q).Valid() = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestValidateTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    State
		to      State
		wantErr bool
	}{
		// Valid transitions
		{"thinking to shaping", Thinking, Shaping, false},
		{"thinking to building (skip shaping)", Thinking, Building, false},
		{"shaping to building", Shaping, Building, false},
		{"building to shipped", Building, Shipped, false},

		// Invalid transitions
		{"thinking to shipped", Thinking, Shipped, true},
		{"shaping to thinking", Shaping, Thinking, true},
		{"shaping to shipped", Shaping, Shipped, true},
		{"shaping to shaping", Shaping, Shaping, true},
		{"building to thinking", Building, Thinking, true},
		{"building to shaping", Building, Shaping, true},
		{"shipped to thinking", Shipped, Thinking, true},
		{"shipped to shaping", Shipped, Shaping, true},
		{"shipped to building", Shipped, Building, true},
		{"thinking to thinking", Thinking, Thinking, true},
		{"building to building", Building, Building, true},
		{"shipped to shipped", Shipped, Shipped, true},

		// Invalid states
		{"invalid from state", State("invalid"), Building, true},
		{"invalid to state", Thinking, State("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransition(%q, %q) error = %v, wantErr %v", tt.from, tt.to, err, tt.wantErr)
			}
		})
	}
}

func TestNextValidActions(t *testing.T) {
	tests := []struct {
		state State
		want  []string
	}{
		{Thinking, []string{"accept", "accept --skip-shaping", "reject", "reset"}},
		{Shaping, []string{"shape", "approve", "revise", "reset"}},
		{Building, []string{"ship", "reset"}},
		{Shipped, []string{"reset"}},
		{State("invalid"), nil},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := NextValidActions(tt.state)
			if len(got) != len(tt.want) {
				t.Errorf("NextValidActions(%q) = %v, want %v", tt.state, got, tt.want)
				return
			}
			for i, action := range got {
				if action != tt.want[i] {
					t.Errorf("NextValidActions(%q)[%d] = %q, want %q", tt.state, i, action, tt.want[i])
				}
			}
		})
	}
}
