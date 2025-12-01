ALTER TABLE links
    DROP CONSTRAINT links_initial_link_key;

CREATE UNIQUE INDEX idx_links_initial_link_unique ON links(initial_link);