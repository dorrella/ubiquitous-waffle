/* Customer Table
 * 
 * Used to track customer info. it makes no assumptions
 * on what is a valid name,as those assumptions can be
 * wrong.
 *
 * Uses empty strings instead of null for qol.
*/


BEGIN;

CREATE TABLE IF NOT EXISTS customers (
       id SERIAL PRIMARY KEY,
       name_pref TEXT NOT NULL,
       name_first TEXT NOT NULL,
       name_middle TEXT NOT NULL,
       name_last TEXT NOT NULL,
       name_suffix TEXT NOT NULL,
       email TEXT NOT NULL UNIQUE,
       phone_number VARCHAR(16) NOT NULL, --longest phone number is around 15 chars
       deleted BOOLEAN NOT NULL,  --customer has been marked for deletion
       created_by INT NOT NULL,
       created_at TIMESTAMP NOT NULL,
       updated_by INT NOT NULL,
       updated_at TIMESTAMP NOT NULL
);

COMMIT;
