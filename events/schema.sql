CREATE TABLE events (
     id TEXT,                                 -- Unique identifier for the event
     source TEXT NOT NULL,                    -- URI reference that identifies the context in which the event happened
     type TEXT NOT NULL,                      -- The type of event
     specversion TEXT NOT NULL,               -- The version of the CloudEvents specification
     datacontenttype TEXT,                    -- Content type of the event data
     data JSONB NOT NULL,                     -- The actual event payload stored as JSONB
     time TIMESTAMPTZ NOT NULL,               -- The timestamp of when the event occurred (ISO 8601)
     subject TEXT,                            -- Optional identifier or description of the event subject
     PRIMARY KEY (id, time)
);
