package parameter_store

import "testing"

func TestParameter(t *testing.T) {
	type args struct {
		secretsStore KeyClient
		key          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "#1: Single Key - Correct Key",
			args: args{
				secretsStore: KeyClient{
					KeyValFile: KeyValFile{
						KeyVals: []struct {
							Key        string                 "json:\"key\""
							Parameters map[string]interface{} "json:\"value\""
						}{
							{
								Key: "dbuser/staging/admin/third_party_api_delete",
							},
						},
					},
				},
				key: "dbuser/staging/admin/third_party_api_delete",
			},
			wantErr: false,
		},
		{
			name: "#2: Single Key - Incorrect Key",
			args: args{
				secretsStore: KeyClient{
					KeyValFile: KeyValFile{
						KeyVals: []struct {
							Key        string                 "json:\"key\""
							Parameters map[string]interface{} "json:\"value\""
						}{
							{
								Key: "dbuser/staging/admin/third_party_api_delete",
							},
						},
					},
				},
				key: "dbuser/staging/admin/third_party_api_deletee",
			},
			wantErr: true,
		},
		{
			name: "#3: Empty",
			args: args{
				secretsStore: KeyClient{
					KeyValFile: KeyValFile{
						KeyVals: nil,
					},
				},
				key: "dbuser/staging/admin/third_party_api_deletee",
			},
			wantErr: true,
		},
		{
			name: "#4: Multiple Keys - Correct Key",
			args: args{
				secretsStore: KeyClient{
					KeyValFile: KeyValFile{
						KeyVals: []struct {
							Key        string                 "json:\"key\""
							Parameters map[string]interface{} "json:\"value\""
						}{
							{
								Key:        "dbuser/staging/admin/first_party_api_delete",
								Parameters: map[string]interface{}{},
							},
							{
								Key:        "dbuser/staging/admin/third_party_api_delete",
								Parameters: map[string]interface{}{},
							},
							{
								Key:        "dbuser/staging/admin/fifth_party_api_delete",
								Parameters: map[string]interface{}{},
							},
						},
					},
				},
				key: "dbuser/staging/admin/third_party_api_delete",
			},
			wantErr: false,
		},
		{
			name: "#5: Multiple Keys - Not found",
			args: args{
				secretsStore: KeyClient{
					KeyValFile: KeyValFile{
						KeyVals: []struct {
							Key        string                 "json:\"key\""
							Parameters map[string]interface{} "json:\"value\""
						}{
							{
								Key:        "dbuser/staging/admin/first_party_api_delete",
								Parameters: map[string]interface{}{},
							},
							{
								Key:        "dbuser/staging/admin/third_party_api_delete",
								Parameters: map[string]interface{}{},
							},
							{
								Key:        "dbuser/staging/admin/fifth_party_api_delete",
								Parameters: map[string]interface{}{},
							},
						},
					},
				},
				key: "dbuser/staging/admin/tenth_party_api_delete",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.args.secretsStore.Parameter(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
