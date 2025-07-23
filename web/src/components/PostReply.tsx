import { useMutation, useQueryClient } from '@tanstack/react-query'

interface IPostReply {
  id: number
  post_id: number
  pub_date: string
  author: string
  content: string
  content_html: string
}

interface PostReplyProps {
  reply: IPostReply
  onDelete: (replyId: number) => void
}

async function deleteReply(replyId: number): Promise<void> {
  const response = await fetch(`/api/replies/${replyId}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to delete reply')
  }
}

export function PostReply({ reply, onDelete }: PostReplyProps) {
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: (replyId: number) => deleteReply(replyId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['replies', reply.post_id] })
      onDelete(reply.id)
    },
  })

  const handleDelete = () => {
    if (window.confirm('Delete this reply?')) {
      deleteMutation.mutate(reply.id)
    }
  }

  return (
    <div 
      className="ml-6 mt-3 p-3 rounded-lg border-l-2 group transition-colors"
      style={{
        backgroundColor: 'rgb(var(--bg-secondary))',
        borderLeftColor: 'rgb(var(--border))',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.backgroundColor = 'rgb(var(--bg-primary))'
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.backgroundColor = 'rgb(var(--bg-secondary))'
      }}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <span 
              className="text-sm font-medium"
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              {reply.author}
            </span>
            <span 
              className="text-xs"
              style={{ color: 'rgb(var(--text-secondary) / 0.7)' }}
            >
              {new Date(reply.pub_date).toLocaleString()}
            </span>
          </div>
          <div 
            className="text-sm prose prose-sm max-w-none"
            style={{ color: 'rgb(var(--text-primary))' }}
            dangerouslySetInnerHTML={{ __html: reply.content_html }}
          />
        </div>
        <button
          onClick={handleDelete}
          disabled={deleteMutation.isPending}
          className="ml-2 p-1 opacity-0 group-hover:opacity-100 transition-all duration-200"
          style={{ 
            color: 'rgb(var(--text-secondary) / 0.6)',
          } as React.CSSProperties}
          onMouseEnter={(e) => {
            e.currentTarget.style.color = '#ef4444' // red-500
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.color = 'rgb(var(--text-secondary) / 0.6)'
          }}
          title="Delete reply"
        >
          {deleteMutation.isPending ? (
            <div 
              className="w-4 h-4 border-2 border-t-transparent rounded-full animate-spin"
              style={{ borderColor: 'rgb(var(--text-secondary) / 0.3)' }}
            ></div>
          ) : (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          )}
        </button>
      </div>
    </div>
  )
}
