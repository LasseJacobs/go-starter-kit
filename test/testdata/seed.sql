
DELETE FROM content.stories WHERE author = '$SEEDER';

INSERT INTO content.stories(id, title, author, votes, url)
    VALUES ('aWfT3ss', 'The Go Starter Project', '$SEEDER', 12, 'www.example.com');

INSERT INTO content.stories(id, title, author, votes, url)
    VALUES ('p0jHXCy', 'GoLang 1.18 release, GENERICS!', '$SEEDER', 401, 'www.example.com');

INSERT INTO content.stories(id, title, author, votes, url)
    VALUES ('nh74RTo', 'Why Go does not need generics', '$SEEDER', 812, 'www.example.com');