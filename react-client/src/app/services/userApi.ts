import {api} from "./api";
import type {User} from "../types";

export const userApi = api.injectEndpoints({
    endpoints: (builder) => ({
        login: builder.mutation<
            {token: string},
            {email: string; password: string}
        >({
            query: (userData) => ({
                url: '/auth/login',
                method: "POST",
                body: userData
            })
        }),
        register: builder.mutation<
            {email: string; password: string; name:string},
            {email: string; password: string; name:string}
        >({
            query: (userData) => ({
                url: '/auth/register',
                method: "POST",
                body: userData
            })
        }),
        current: builder.query<User, void>({
            query: () => ({
                url: '/api/current',
                method: "GET"
            })
        }),
        getUserByID: builder.query<User, string>({
            query: (id) => ({
                url: `/api/users/${id}`,
                method: "GET"
            })
        }),
        updateUser: builder.mutation<User, {userData: FormData, id: string}>({
            query: ({userData, id}) => ({
                url: `/api/users/${id}`,
                method: "PUT",
                body: userData
            })
        })
    })
})

export const {
    useRegisterMutation,
    useLoginMutation,
    useLazyCurrentQuery,
    useCurrentQuery,
    useGetUserByIDQuery,
    useLazyGetUserByIDQuery,
    useUpdateUserMutation,
} = userApi;

export const {
    endpoints: {login, register, current, getUserByID, updateUser}
} = userApi