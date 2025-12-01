-- noinspection SqlResolve
ALTER TABLE links
    ADD CONSTRAINT links_initial_link_key
    UNIQUE USING INDEX idx_links_initial_link_unique;