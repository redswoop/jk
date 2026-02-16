package main

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantAll     bool
		wantQuiet   bool
		wantGUI     bool
		wantPorts   bool
		wantTree    bool
		wantFilter  string
		wantHelp    bool
		wantVersion bool
	}{
		{
			"no args",
			[]string{},
			false, false, false, false, false, "", false, false,
		},
		{
			"all flags short",
			[]string{"-a", "-q", "-g", "-p", "-t"},
			true, true, true, true, true, "", false, false,
		},
		{
			"all flags long",
			[]string{"--all", "--quiet", "--gui-apps", "--ports", "--tree"},
			true, true, true, true, true, "", false, false,
		},
		{
			"filter only",
			[]string{"node"},
			false, false, false, false, false, "node", false, false,
		},
		{
			"filter with flags interleaved",
			[]string{"-p", "node", "-q"},
			false, true, false, true, false, "node", false, false,
		},
		{
			"filter is lowercased",
			[]string{"Node"},
			false, false, false, false, false, "node", false, false,
		},
		{
			"help short",
			[]string{"-h"},
			false, false, false, false, false, "", true, false,
		},
		{
			"help long",
			[]string{"--help"},
			false, false, false, false, false, "", true, false,
		},
		{
			"first non-flag wins as filter",
			[]string{"node", "python"},
			false, false, false, false, false, "node", false, false,
		},
		{
			"unknown flags ignored",
			[]string{"--unknown", "-x"},
			false, false, false, false, false, "", false, false,
		},
		{
			"version short",
			[]string{"-V"},
			false, false, false, false, false, "", false, true,
		},
		{
			"version long",
			[]string{"--version"},
			false, false, false, false, false, "", false, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all, quiet, gui, ports, tree, filter, help, ver := parseArgs(tt.args)
			if all != tt.wantAll {
				t.Errorf("showAll = %v, want %v", all, tt.wantAll)
			}
			if quiet != tt.wantQuiet {
				t.Errorf("quiet = %v, want %v", quiet, tt.wantQuiet)
			}
			if gui != tt.wantGUI {
				t.Errorf("showGUI = %v, want %v", gui, tt.wantGUI)
			}
			if ports != tt.wantPorts {
				t.Errorf("showPorts = %v, want %v", ports, tt.wantPorts)
			}
			if tree != tt.wantTree {
				t.Errorf("showTree = %v, want %v", tree, tt.wantTree)
			}
			if filter != tt.wantFilter {
				t.Errorf("filterTerm = %q, want %q", filter, tt.wantFilter)
			}
			if help != tt.wantHelp {
				t.Errorf("showHelp = %v, want %v", help, tt.wantHelp)
			}
			if ver != tt.wantVersion {
				t.Errorf("showVersion = %v, want %v", ver, tt.wantVersion)
			}
		})
	}
}
