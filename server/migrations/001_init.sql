CREATE TABLE IF NOT EXISTS schema_migrations (
  schema_migration_id INTEGER PRIMARY KEY AUTOINCREMENT,
  schema_migration_filename TEXT NOT NULL UNIQUE,
  schema_migration_applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS settings (
  setting_key TEXT UNIQUE NOT NULL,
  setting_json BLOB NOT NULL,
  setting_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  setting_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
  user_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_username TEXT UNIQUE NOT NULL,
  user_password TEXT NOT NULL,
  user_email TEXT,
  user_role TEXT DEFAULT 'user',
  user_is_active BOOLEAN DEFAULT TRUE,
  user_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  user_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  user_last_login DATETIME
);


CREATE TABLE IF NOT EXISTS revoked_tokens (
  user_id INTEGER NOT NULL,
  token TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS personal_infos (
  personal_info_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  personal_info_first_name TEXT NOT NULL,
  personal_info_last_name TEXT NOT NULL,
  personal_info_professional_title TEXT,
  personal_info_bio TEXT,
  personal_info_location TEXT,
  personal_info_resume_url TEXT,
  personal_info_website_url TEXT,
  personal_info_linkedin_url TEXT,
  personal_info_github_url TEXT,
  personal_info_x_url TEXT,
  personal_info_date_of_birth DATE,
  personal_info_phone_number TEXT,
  personal_info_interests TEXT,
  personal_info_profile_picture TEXT,
  personal_info_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  personal_info_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS projects (
  project_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  project_title TEXT NOT NULL,
  project_description TEXT,
  project_short_description TEXT,
  project_technologies TEXT,
  project_github_url TEXT,
  project_image_url TEXT,
  project_status TEXT DEFAULT 'active',
  project_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  project_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(project_title, user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS skills (
  skill_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  skill_name TEXT NOT NULL,
  skill_level INTEGER CHECK(skill_level BETWEEN 1 AND 5),
  skill_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  skill_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(skill_name, user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS experiences (
  experience_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  experience_job_title TEXT NOT NULL,
  experience_company_name TEXT NOT NULL,
  experience_start_date DATE,
  experience_end_date DATE,
  experience_description TEXT,
  experience_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  experience_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(experience_job_title, experience_company_name, user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS educations (
  education_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  education_degree TEXT NOT NULL,
  education_institution TEXT NOT NULL,
  education_start_date DATE,
  education_end_date DATE,
  education_description TEXT,
  education_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  education_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(education_degree, education_institution, user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS technologies (
  technology_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  technology_name TEXT NOT NULL,
  technology_icon_url TEXT NOT NULL,
  technology_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  technology_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(technology_name, user_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
