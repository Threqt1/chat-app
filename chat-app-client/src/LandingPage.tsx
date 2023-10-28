import { Button, Center, HStack, Input, VStack, useToast } from "@chakra-ui/react"
import { useState } from "react"
import { useLoaderData, useNavigate } from "react-router-dom"
import { ErrorInvalidCredentials, ErrorUsernameDoesntExist, ErrorUsernameNotUnique, Login, Register } from "./APIWrapper"
import { WebsocketLoaderResult } from "./main"

const ErrorUnknown = "An unknown error occurred."

export const LandingPage = () => {
    const { ws } = useLoaderData() as WebsocketLoaderResult
    const navigate = useNavigate();

    const [username, setUsername] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [loading, setLoading] = useState<boolean>(false)
    const [invalid, setInvalid] = useState<boolean>(false)
    const errorToast = useToast();

    const invalidInput = (type: string, error: string) => {
        errorToast({
            title: `Failed to ${type}.`,
            description: error,
            status: "error",
            duration: 2000,
            isClosable: true
        })
        setLoading(false)
        setInvalid(true)
    }

    const onSubmit = async (username: string) => {
        await ws.connect(username, 3000)
        ws.setUsername(username)
        navigate("/")
    }

    const onClickLogin = async () => {
        setLoading(true)
        setInvalid(false)
        if (username.length === 0 || password.length === 0) {
            return invalidInput("login", "Username or password left empty.")
        }
        let errorMsg = ErrorUnknown
        try {
            const data = await Login(username, password)
            if (data.code !== 200) {
                switch (data.message) {
                    case ErrorUsernameDoesntExist:
                        errorMsg = "Username doesn't exist."
                        break;
                    case ErrorInvalidCredentials:
                        errorMsg = "Invalid credentials."
                        break;
                }
                return invalidInput("login", errorMsg)
            }
            if (!data.data) {
                return invalidInput("login", errorMsg)
            }
            await onSubmit(data.data.user.username)
        } catch (err: unknown) {
            console.error(err)
            return invalidInput("login", errorMsg)
        }
    }

    const onClickRegister = async () => {
        setLoading(true)
        setInvalid(false)
        if (username.length === 0 || password.length === 0) {
            return invalidInput("register", "Username or password left empty.")
        }
        let errorMsg = ErrorUnknown
        try {
            const data = await Register(username, password)
            if (data.code !== 200) {
                switch (data.message) {
                    case ErrorUsernameNotUnique:
                        errorMsg = "The username is not unique."
                        break;
                }
                return invalidInput("register", errorMsg)
            }
            if (!data.data) {
                return invalidInput("register", errorMsg)
            }
            await onSubmit(data.data.user.username)
        } catch (err: unknown) {
            console.error(err)
            return invalidInput("register", errorMsg)
        }
    }

    return <Center w="100vw" h="100vh">
        <VStack spacing={"15px"}>
            <Input isInvalid={invalid} boxShadow={"md"} placeholder="Username" errorBorderColor="red.400" textAlign={"center"} value={username} onChange={(event) => {
                if (event.target.value === "" || (/^[a-zA-Z0-9_]*$/).test(event.target.value)) setUsername(event.target.value.trim())
            }} />
            <Input isInvalid={invalid} boxShadow={"md"} placeholder="Password" errorBorderColor="red.400" textAlign={"center"} type={"password"} onChange={(event) => setPassword(event.target.value.trim())} />
            <HStack>
                <Button isLoading={loading} boxShadow={"md"} size={"md"} onClick={onClickLogin}>Login</Button>
                <Button isLoading={loading} boxShadow={"md"} size={"md"} onClick={onClickRegister}>Register</Button>
            </HStack>
        </VStack>
    </Center>
}