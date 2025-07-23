import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import { markTopicRead } from '../../../../utils/readStatus'
import { PostWithReplies } from '../../../../components/PostWithReplies'

interface ITopic {
  id: number
  board_id: number
  pub_date: string
  title: string
  status: string
  author: string
  last_post_id: number
  post_count: number
}

interface IPost {
  id: number
  topic_id: number
  pub_date: string
  author: string
  content: string
  content_html: string
}

async function fetchTopic(topicId: string): Promise<ITopic> {
  const response = await fetch(`/api/topic/${topicId}`)
  if (!response.ok) {
    throw new Error('Failed to fetch topic')
  }
  return response.json()
}

async function fetchPosts(topicId: string): Promise<IPost[]> {
  const response = await fetch(`/api/posts/topic/${topicId}`)
  if (!response.ok) {
    throw new Error('Failed to fetch posts')
  }
  const data = await response.json()
  return data.posts
}

async function createPost(topicId: number, author: string, content: string): Promise<IPost> {
  const response = await fetch('/api/post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      topic_id: topicId,
      author,
      content,
    }),
  })
  if (!response.ok) {
    throw new Error('Failed to create post')
  }
  return response.json()
}

function Topic() {
  const { slug, topicId } = Route.useParams()
  const [newPostAuthor, setNewPostAuthor] = useState('')
  const [newPostContent, setNewPostContent] = useState('')
  const queryClient = useQueryClient()

  const topicQuery = useQuery({
    queryKey: ['topic', topicId],
    queryFn: () => fetchTopic(topicId),
    enabled: !!topicId,
  })

  const postsQuery = useQuery({
    queryKey: ['posts', topicId],
    queryFn: () => fetchPosts(topicId),
    enabled: !!topicId,
  })

  const createPostMutation = useMutation({
    mutationFn: ({ author, content }: { author: string; content: string }) =>
      createPost(Number(topicId), author, content),
    onSuccess: (newPost) => {
      queryClient.invalidateQueries({ queryKey: ['posts', topicId] })
      queryClient.invalidateQueries({ queryKey: ['topic', topicId] })
      setNewPostContent('')
      if (topicQuery.data) {
        markTopicRead(Number(topicId), newPost.id)
      }
    },
  })

  useEffect(() => {
    if (topicQuery.data && postsQuery.data && postsQuery.data.length > 0) {
      const lastPostId = Math.max(...postsQuery.data.map((p: IPost) => p.id))
      markTopicRead(Number(topicId), lastPostId)
    }
  }, [topicQuery.data, postsQuery.data, topicId])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (newPostContent.trim()) {
      createPostMutation.mutate({
        author: newPostAuthor || 'Anonymous',
        content: newPostContent,
      })
    }
  }

  if (topicQuery.isLoading || postsQuery.isLoading) {
    return <div className="p-6">Loading...</div>
  }

  if (topicQuery.error || postsQuery.error) {
    return (
      <div className="p-6 text-red-600">
        Error loading topic: {topicQuery.error?.message || postsQuery.error?.message}
      </div>
    )
  }

  const topic = topicQuery.data
  const posts = postsQuery.data || []

  return (
    <div className="container mx-auto px-4 py-6">
      <div className="mb-4">
        <Link 
          to="/board/$slug" 
          params={{ slug }} 
          className="text-blue-600 hover:text-blue-800"
        >
          ← Back to /{slug}/
        </Link>
      </div>

      {topic && (
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">{topic.title}</h1>
          <div className="text-gray-600">
            Started by {topic.author} • {new Date(topic.pub_date).toLocaleString()} • {topic.post_count} posts
          </div>
        </div>
      )}

      <div className="space-y-6 mb-8">
        {posts.map((post) => (
          <PostWithReplies key={post.id} post={post} />
        ))}
      </div>

      <form onSubmit={handleSubmit} className="border border-gray-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold mb-4">Reply to this topic</h3>
        
        <div className="mb-4">
          <label htmlFor="author" className="block text-sm font-medium text-gray-700 mb-2">
            Name (optional)
          </label>
          <input
            type="text"
            id="author"
            value={newPostAuthor}
            onChange={(e) => setNewPostAuthor(e.target.value)}
            placeholder="Anonymous"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
            Message (supports Markdown)
          </label>
          <textarea
            id="content"
            value={newPostContent}
            onChange={(e) => setNewPostContent(e.target.value)}
            rows={6}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Write your reply here..."
          />
        </div>

        <button
          type="submit"
          disabled={!newPostContent.trim() || createPostMutation.isPending}
          className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {createPostMutation.isPending ? 'Posting...' : 'Post Reply'}
        </button>
      </form>
    </div>
  )
}

export const Route = createFileRoute('/board_/$slug/topic/$topicId')({
  component: Topic,
})
