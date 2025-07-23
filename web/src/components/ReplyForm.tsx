import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

interface IPostReply {
  id: number
  post_id: number
  pub_date: string
  author: string
  content: string
  content_html: string
}

interface ReplyFormProps {
  postId: number
  onCancel: () => void
  onSuccess?: (reply: IPostReply) => void
}

async function createReply(postId: number, author: string, content: string): Promise<IPostReply> {
  const response = await fetch(`/api/posts/${postId}/replies`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      author,
      content,
    }),
  })
  if (!response.ok) {
    throw new Error('Failed to create reply')
  }
  return response.json()
}

export function ReplyForm({ postId, onCancel, onSuccess }: ReplyFormProps) {
  const [author, setAuthor] = useState('')
  const [content, setContent] = useState('')
  const queryClient = useQueryClient()

  const createMutation = useMutation({
    mutationFn: ({ author, content }: { author: string; content: string }) =>
      createReply(postId, author, content),
    onSuccess: (newReply) => {
      queryClient.invalidateQueries({ queryKey: ['replies', postId] })
      setAuthor('')
      setContent('')
      onSuccess?.(newReply)
      onCancel()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!content.trim()) return
    
    createMutation.mutate({ 
      author: author.trim() || 'Anonymous', 
      content: content.trim() 
    })
  }

  return (
    <form 
      onSubmit={handleSubmit} 
      className="ml-6 mt-3 p-3 rounded-lg border"
      style={{
        backgroundColor: 'rgb(var(--accent) / 0.1)',
        borderColor: 'rgb(var(--accent) / 0.3)',
      }}
    >
      <div className="space-y-3">
        <input
          type="text"
          placeholder="Your name (optional)"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          className="w-full px-3 py-2 text-sm border rounded focus:outline-none focus:ring-2 transition-all"
          style={{
            backgroundColor: 'rgb(var(--bg-primary))',
            borderColor: 'rgb(var(--border))',
            color: 'rgb(var(--text-primary))',
            '--tw-ring-color': 'rgb(var(--accent))',
          } as React.CSSProperties}
        />
        <textarea
          placeholder="Write your reply... (Markdown supported)"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="w-full px-3 py-2 text-sm border rounded focus:outline-none focus:ring-2 resize-none transition-all"
          style={{
            backgroundColor: 'rgb(var(--bg-primary))',
            borderColor: 'rgb(var(--border))',
            color: 'rgb(var(--text-primary))',
            '--tw-ring-color': 'rgb(var(--accent))',
          } as React.CSSProperties}
          rows={3}
          required
        />
        <div className="flex gap-2">
          <button
            type="submit"
            disabled={createMutation.isPending || !content.trim()}
            className="btn-primary text-sm disabled:cursor-not-allowed"
          >
            {createMutation.isPending ? 'Posting...' : 'Reply'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="btn-secondary text-sm"
          >
            Cancel
          </button>
        </div>
      </div>
    </form>
  )
}
