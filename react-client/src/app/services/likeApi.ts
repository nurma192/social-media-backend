import {api} from "./api";
import type {Like} from "../types";

export const likeApi = api.injectEndpoints({
    endpoints: (builder) => ({
        likePost: builder.mutation<Like, { postID: string }>({
            query: (body) => ({
                url: `/likes`,
                method: 'POST',
                body,
            })
        }),
        unLikePost: builder.mutation<void, string>({
            query: (postID) => ({
                url: `/like/${postID}`,
                method: 'DELETE',
            })
        })
    })
})

export const {
    useLikePostMutation,
    useUnLikePostMutation,
} = likeApi

export const {
    endpoints: {likePost, unLikePost}
} = likeApi