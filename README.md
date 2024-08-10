![Alt text](english-terminal.png){ width=60% }

# english-terminal

openAI powered terminal. write commands in english and get the output in terminal.

example input output

```
english> is port 1234 taken
COMMAND     PID      USER   FD   TYPE            DEVICE SIZE/OFF NODE NAME
LM\x20Stu 38858 xxxxxxxxx   93u  IPv6 0x6f227a26a957be5      0t0  TCP *:search-agent (LISTEN)


english> list running containers
35340e23ed36   mysql:8.0               "docker-entrypoint.s…"   3 days ago   Up 3 days   0.0.0.0:3306->3306/tcp, 33060/tcp   mysql_db
```


# Installation:
1) download english executable and copy it in a path and reference it in your PATH variable.
2) OPENAI_API_KEY should be set in environment variables.
3) Usage :  type english in your terminal. You should see a english> prompt.


# Modes: 
? (Answer only)
english> ?<command in english>
only displays the command, doesnt run.

Example: 
```
english>? is port 1234 taken
lsof -i :1234
```


!(see command before running)
english> !<command in english>
Displays the command and waits for Enter key and then runs it.
Example:
```
english>! delete folder lib
:> rm -rf lib
Press enter to continue...
```




