from subprocess import Popen, PIPE
import os
import time
import asyncio
import locale
import sys
from asyncio.subprocess import PIPE
from contextlib import closing

serverArgs = ["./oauth2local",
              "--verbose",
              "--config",
              "test-oauth2-config.yml",
              "serve"]
callbackArgs = ["./oauth2local",
                "--verbose",
                "--config",
                "test-oauth2-config.yml",
                "callback",
                "loc-auth://callback?code=cryptocode"]
tokenArgs = ["./oauth2local",
             "--verbose",
             "--config",
             "test-oauth2-config.yml",
             "token"]


async def readline_and_kill(args, sf, cf, tf):
    # start child process
    server = await asyncio.create_subprocess_exec(
        args[0], *args[1:], stdout=sf, stderr=sf)
    time.sleep(3)
    callbackCmd = Popen(callbackArgs,
                        stdout=cf,
                        stderr=cf)
    time.sleep(3)
    callbackCmd = Popen(callbackArgs,
                        stdout=tf,
                        stderr=tf)
    # read line (sequence of bytes ending with b'\n') asynchronously
    # async for line in process.stdout:
    #     print("got line:", line.decode(locale.getpreferredencoding(False)))
    # break
    # process.kill()
    print("Awaiting")
    return await server.wait()  # wait for the child process to exit


if sys.platform == "win32":
    loop = asyncio.ProactorEventLoop()
    asyncio.set_event_loop(loop)
else:
    loop = asyncio.get_event_loop()

with closing(loop):
    sys.exit(loop.run_until_complete(
        readline_and_kill(
            serverArgs,
            open("server.log", "w"),
            open("callback.log", "w"),
            open("token.log", "w"),
        )))
