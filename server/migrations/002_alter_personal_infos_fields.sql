-- Migration: Rename personal_info_bio to personal_info_intro and add personal_info_about_me

-- Rename column personal_info_bio to personal_info_intro
ALTER TABLE personal_infos RENAME COLUMN personal_info_bio TO personal_info_intro;

-- Add new field personal_info_about_me
ALTER TABLE personal_infos ADD COLUMN personal_info_about_me TEXT;
