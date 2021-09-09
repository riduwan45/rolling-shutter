CREATE SCHEMA IF NOT EXISTS keyper;
CREATE TABLE IF NOT EXISTS keyper.decryption_trigger (
       epoch_id bigint PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS keyper.decryption_key_share (
       epoch_id bigint,
       keyper_index bigint,
       decryption_key_share bytea,
       PRIMARY KEY (epoch_id, keyper_index)
);
CREATE TABLE IF NOT EXISTS keyper.decryption_key (
       epoch_id bigint PRIMARY KEY,
       keyper_index bigint,
       decryption_key bytea
);