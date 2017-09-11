package config


type LogConfig struct {

	Active 						bool 			`default:"true"`

	Access struct {
		AccessLogFilePath      	string 			`json:"access_log_filepath,omitempty" yaml:"access_log_filepath,omitempty"`
		AccessLogFileExtension 	string 			`json:"access_log_file_extension,omitempty" yaml:"access_log_file_extension,omitempty"`
		AccessLogMaxSize       	int    			`json:"access_log_max_size,omitempty" yaml:"access_log_max_size,omitempty"`
		AccessLogMaxBackups    	int    			`json:"access_log_max_backups,omitempty" yaml:"access_log_max_backups,omitempty"`
		AccessLogMaxAge        	int    			`json:"access_log_max_age,omitempty" yaml:"access_log_max_age,omitempty"`
	}

	Error struct {
		ErrorLogFilePath       	string 			`json:"error_log_filepath,omitempty" yaml:"error_log_filepath,omitempty"`
		ErrorLogFileExtension  	string 			`json:"error_log_file_extension,omitempty" yaml:"error_log_file_extension,omitempty"`
		ErrorLogMaxSize        	int    			`json:"error_log_max_size,omitempty" yaml:"error_log_max_size,omitempty"`
		ErrorLogMaxBackups     	int    			`json:"error_log_max_backups,omitempty" yaml:"error_log_max_backups,omitempty"`
		ErrorLogMaxAge         	int    			`json:"error_log_max_age,omitempty" yaml:"error_log_max_age,omitempty"`
	}

}