import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation } from '@tanstack/react-query'
import { useState, useEffect } from 'react'

interface IBoard {
  id: number
  slug: string
  description: string
}

async function fetchBoard(slug: string): Promise<IBoard> {
  const response = await fetch(`/api/board/${slug}`)
  if (!response.ok) {
    throw new Error('Failed to fetch board')
  }
  return response.json()
}

async function createTopic(boardId: number, title: string, author: string): Promise<any> {
  const response = await fetch('/api/topic', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      board_id: boardId,
      title,
      author,
    }),
  })
  if (!response.ok) {
    throw new Error('Failed to create topic')
  }
  return response.json()
}

function NewTopic() {
  console.log('🎯 NewTopic component loaded')
  const { slug } = Route.useParams()
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [author, setAuthor] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const { data: board, isLoading } = useQuery({
    queryKey: ['board', slug],
    queryFn: () => fetchBoard(slug),
  })

  const createTopicMutation = useMutation({
    mutationFn: ({ title, author }: { title: string; author: string }) =>
      createTopic(board!.id, title, author),
    onSuccess: (topic) => {
      navigate({
        to: '/board/$slug/topic/$topicId',
        params: { slug, topicId: topic.id.toString() },
      })
    },
    onError: (error) => {
      console.error('Failed to create topic:', error)
      setIsSubmitting(false)
    },
  })

  const closeModal = () => {
    navigate({ to: '/board/$slug', params: { slug } })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim() || !author.trim() || !board) return
    
    setIsSubmitting(true)
    createTopicMutation.mutate({ title, author })
  }

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      closeModal()
    }
  }

  // Handle escape key and prevent body scroll
  useEffect(() => {
    const handleEscapeKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        closeModal()
      }
    }

    // Prevent body scroll when modal is open
    document.body.style.overflow = 'hidden'
    document.addEventListener('keydown', handleEscapeKey)
    
    return () => {
      // Re-enable body scroll when modal closes
      document.body.style.overflow = 'unset'
      document.removeEventListener('keydown', handleEscapeKey)
    }
  }, [])

  if (isLoading) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="themed-bg-secondary rounded-lg p-6">
          <div className="text-center text-primary">Loading...</div>
        </div>
      </div>
    )
  }

  return (
    <div 
      className="fixed inset-0 modal-backdrop flex items-center justify-center z-50 p-4"
      onClick={handleBackdropClick}
      style={{ top: 0, left: 0, right: 0, bottom: 0, position: 'fixed' }}
    >
      <div className="themed-bg-secondary rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        {/* Modal Header */}
        <div className="flex items-center justify-between p-6 border-b themed-border">
          <div>
            <h1 className="text-2xl font-bold text-primary">Create New Topic</h1>
            <p className="text-secondary">in /{board?.slug}/</p>
          </div>
          <button
            onClick={closeModal}
            className="modal-close-btn"
            aria-label="Close modal"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Modal Body */}
        <div className="p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-primary">
                Topic Title
              </label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="form-input"
                placeholder="Enter topic title..."
                required
                autoFocus
              />
            </div>
            
            <div>
              <label htmlFor="author" className="block text-sm font-medium text-primary">
                Name (use ##password for tripcode)
              </label>
              <input
                type="text"
                id="author"
                value={author}
                onChange={(e) => setAuthor(e.target.value)}
                className="form-input"
                placeholder="Anonymous"
                required
              />
              <p className="mt-1 text-xs text-secondary">
                Use ##yourpassword after your name to generate a tripcode for authentication
              </p>
            </div>

            <div className="flex justify-end space-x-4 pt-4">
              <button
                type="button"
                onClick={closeModal}
                className="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting || !title.trim() || !author.trim()}
                className="btn-primary disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Creating...' : 'Create Topic'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export const Route = createFileRoute('/board/$slug/new-topic')({
  component: NewTopic
})
