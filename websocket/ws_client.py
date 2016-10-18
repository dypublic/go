#!/usr/bin/env python3
import asyncio
import aiohttp


async def main():
    session = aiohttp.ClientSession()
    async with session.ws_connect('ws://127.0.0.1:8080/echo') as ws:

        async for msg in ws:
            if msg.type == aiohttp.WSMsgType.TEXT:
                if msg.data == 'close cmd':
                    await ws.close()
                    break
                else:
                    ws.send_str(msg.data + '/answer')
            elif msg.type == aiohttp.WSMsgType.CLOSED:
                break
            elif msg.type == aiohttp.WSMsgType.ERROR:
                break


# async def main(loop):
#     async with aiohttp.ClientSession(loop=loop) as client:
#         html = await fetch(client)
#         print(html)

if __name__ == '__main__':

    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())