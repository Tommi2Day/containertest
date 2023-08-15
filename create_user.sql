create user c##test identified by "Test!Password"
    default tablespace users quota unlimited on users
    temporary tablespace temp
    account unlock;

grant create session to c##test container=all;
grant set container to c##test container=all;