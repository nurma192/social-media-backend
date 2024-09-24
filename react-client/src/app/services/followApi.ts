import {api} from "./api";

export const followApi = api.injectEndpoints({
    endpoints: (builder) => ({
        followUser: builder.mutation<void, { followingID: string }>({
            query: (body) => ({
                url: `/follow`,
                method: 'POST',
                body
            })
        }),
        unFollowUser: builder.mutation<void, string>({
            query: (userID) => ({
                url: `unfollow/${userID}`,
                method: 'DELETE'
            })
        })
    })
})

export const {
    useFollowUserMutation,
    useUnFollowUserMutation,
} = followApi

export const {
    endpoints: {followUser, unFollowUser}
} = followApi