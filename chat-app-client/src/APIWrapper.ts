import axios from "axios"
import { TChannel, TMessage } from "./types"

export const ErrorUsernameDoesntExist = "one of the usernames doesnt exist"
export const ErrorUsernameNotUnique = "username is not unique"
export const ErrorInvalidCredentials = "invalid credentials"
export const ErrorChannelDoesntExist = "channel doesnt exist"

const BASE = "http://localhost:8080"
const REGISTER_ENDPOINT = "/user/register"
const LOGIN_ENDPOINT = "/user/login"
const GET_CHANNELS_ENDPOINT = "/user/channels"
const CREATE_CHANNEL_ENDPOINT = "/channels/create"
const GET_CHANNEL_ENDPOINT = "/channels/get"
const SEARCH_MESSAGES_ENDPOINT = "/channels/search"


interface HTTPResponse<Res> {
    code: number
    message: string
    data?: Res
}

type HTTPError = { response?: { status: number, data: string } }

async function Request<Req, Res>(request: Req, endpoint: string): Promise<HTTPResponse<Res>> {
    try {
        const { data } = await axios.post<Res>(BASE + endpoint, JSON.stringify(request))
        return { code: 200, message: "", data }
    } catch (err: unknown) {
        const typedErr = err as HTTPError
        if (typedErr.response) {
            return {
                code: typedErr.response.status,
                message: typedErr.response.data.toLowerCase().trim()
            }
        }
        throw new Error()
    }
}

interface RegisterRequest {
    username: string
    password: string
}

interface RegisterResponse {
    user: {
        username: string
        timestamp: number
    }
}

export function Register(username: string, password: string) {
    return Request<RegisterRequest, RegisterResponse>({ username, password }, REGISTER_ENDPOINT);
}

interface LoginRequest {
    username: string
    password: string
}

interface LoginResponse {
    user: {
        username: string
        timestamp: number
    }
}

export function Login(username: string, password: string) {
    return Request<LoginRequest, LoginResponse>({ username, password }, LOGIN_ENDPOINT);
}

interface GetChannelsRequest {
    username: string
    start: number
    end: number
}

interface GetChannelsResponse {
    channels: TChannel[]
}

export function GetChannels(username: string, start: number, end: number) {
    return Request<GetChannelsRequest, GetChannelsResponse>({ username, start, end }, GET_CHANNELS_ENDPOINT)
}

interface CreateChannelRequest {
    members: string[]
}

interface CreateChannelResponse {
    channel: TChannel
}

export function CreateChannel(members: string[]) {
    return Request<CreateChannelRequest, CreateChannelResponse>({ members }, CREATE_CHANNEL_ENDPOINT)
}

interface GetChannelRequest {
    channelID: string
}

interface GetChannelResponse {
    channel: TChannel
}

export function GetChannel(channelID: string) {
    return Request<GetChannelRequest, GetChannelResponse>({ channelID }, GET_CHANNEL_ENDPOINT)
}

interface SearchMessagesRequest {
    channelID: string
    start: number
    stop: number
}

interface SearchMessagesResponse {
    messages: TMessage[]
}

export function SearchMessages(channelID: string, start: number, stop: number) {
    return Request<SearchMessagesRequest, SearchMessagesResponse>({ channelID, start, stop }, SEARCH_MESSAGES_ENDPOINT)
}