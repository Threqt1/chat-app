import * as React from 'react'
import { ChakraProvider } from '@chakra-ui/react'
import * as ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider, redirect, Navigate } from "react-router-dom"

import { LandingPage } from './LandingPage'
import SocketConnection from './SocketConnection'
import { MainPage } from './MainPage'
import { ChannelsPage } from './ChannelsPage'
import { MessagesPage } from './MessagesPage'
import { GetChannels, SearchMessages } from './APIWrapper'

const ws = new SocketConnection()
ws.initialize()

export interface WebsocketLoaderResult {
  ws: SocketConnection
}

const router = createBrowserRouter([
  {
    path: "/",
    loader: async () => {
      if (ws.getUsername().length === 0) return redirect("/login")
      return { ws }
    },
    element: <MainPage />,
    children: [
      {
        index: true,
        element: <Navigate to={"/channels"} />
      },
      {
        path: "channels",
        loader: async () => {
          //add load more
          const channels = await GetChannels(ws.getUsername(), 0, 50)
          if (!channels.data) throw new Response("Unable to Load", { status: 404 })
          return { ws, channels: channels.data }
        },
        element: <ChannelsPage />
      },
      {
        path: "message/:channelID",
        loader: async ({ params }) => {
          if (!params.channelID) throw new Response("Unable to Load Params", { status: 404 })
          const messages = await SearchMessages(params.channelID, 0, 50)
          if (!messages.data) throw new Response("Unable to Load", { status: 404 })
          return { ws, channel: params.channelID, messages: messages.data.messages }
        },
        element: <MessagesPage />
      }
    ]
  },
  {
    path: "/login",
    element: <LandingPage />,
    loader: async () => {
      return { ws }
    },
  },
])

const rootElement = document.getElementById('root')!
ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <ChakraProvider>
      <RouterProvider router={router} />
    </ChakraProvider>
  </React.StrictMode>,
)