import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { PostReply } from './PostReply'
import { ReplyForm } from './ReplyForm'

interface IPost {
  id: number
  topic_id: number
  pub_date: string
  author: string
  content: string
  content_html: string
}

interface IPostReply {
  id: number
  post_id: number
  pub_date: string
  author: string
  content: string
  content_html: string
}

interface PostWithRepliesProps {
  post: IPost
}

async function fetchReplies(postId: number): Promise<IPostReply[]> {
  const response = await fetch(`/api/posts/${postId}/replies`)
  if (!response.ok) {
    throw new Error('Failed to fetch replies')
  }
  const data = await response.json()
  return data.replies || []
}

export function PostWithReplies({ post }: PostWithRepliesProps) {
  const [showReplyForm, setShowReplyForm] = useState(false)
  
  const { data: replies = [], isLoading } = useQuery({
    queryKey: ['replies', post.id],
    queryFn: () => fetchReplies(post.id),
  })

  return (
    <div 
      className="border rounded-lg p-6"
      style={{
        borderColor: 'rgb(var(--border))',
        backgroundColor: 'rgb(var(--bg-secondary))',
      }}
    >
      <div className="flex justify-between items-start mb-3">
        <div 
          className="font-semibold"
          style={{ color: 'rgb(var(--accent))' }}
        >
          {post.author}
        </div>
        <div 
          className="text-sm"
          style={{ color: 'rgb(var(--text-secondary))' }}
        >
          {new Date(post.pub_date).toLocaleString()}
        </div>
      </div>
      <div 
        className="prose prose-sm max-w-none mb-4"
        style={{ color: 'rgb(var(--text-primary))' }}
        dangerouslySetInnerHTML={{ __html: post.content_html }}
      />
      
      <div 
        className="flex items-center gap-4 text-sm border-t pt-3"
        style={{
          color: 'rgb(var(--text-secondary))',
          borderColor: 'rgb(var(--border))',
        }}
      >
        <button
          onClick={() => setShowReplyForm(!showReplyForm)}
          className="flex items-center gap-1 transition-colors"
          style={{ color: 'rgb(var(--text-secondary))' }}
          onMouseEnter={(e) => {
            e.currentTarget.style.color = 'rgb(var(--accent))'
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.color = 'rgb(var(--text-secondary))'
          }}
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          Reply {replies.length > 0 && `(${replies.length})`}
        </button>
      </div>

      {showReplyForm && (
        <ReplyForm 
          postId={post.id} 
          onCancel={() => setShowReplyForm(false)}
        />
      )}

      {isLoading && replies.length === 0 && (
        <div 
          className="ml-6 mt-3 text-sm"
          style={{ color: 'rgb(var(--text-secondary))' }}
        >
          Loading replies...
        </div>
      )}

      {replies.length > 0 && (
        <div className="mt-4">
          {replies.map((reply) => (
            <PostReply 
              key={reply.id} 
              reply={reply} 
              onDelete={() => {
                // Reply component will handle the deletion and cache invalidation
              }}
            />
          ))}
        </div>
      )}
    </div>
  )
}
