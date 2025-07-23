import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'

// API types
interface IBoard {
  id: number
  slug: string
  description: string
  topic_count: number
}

async function fetchBoards(): Promise<IBoard[]> {
  const response = await fetch('/api/boards')
  if (!response.ok) {
    throw new Error('Failed to fetch boards')
  }
  const data = await response.json()
  return data.boards
}

function Boards() {
  const { data: boards, isLoading, error } = useQuery({
    queryKey: ['boards'],
    queryFn: fetchBoards,
  })

  if (isLoading) {
    return <div className="text-center py-8">Loading boards...</div>
  }

  if (error) {
    return <div className="text-center py-8 text-red-600">Failed to load boards</div>
  }

  return (
    <div className="space-y-6">
      <div className="rounded-lg shadow-md overflow-hidden" style={{ backgroundColor: 'rgb(var(--bg-secondary))' }}>
        <div className="px-6 py-4 border-b" style={{ 
          backgroundColor: 'rgb(var(--bg-primary))', 
          borderColor: 'rgb(var(--border))' 
        }}>
          <h1 className="text-2xl font-bold" style={{ color: 'rgb(var(--text-primary))' }}>
            Discussion Boards
          </h1>
        </div>
        
        {boards?.length === 0 ? (
          <div className="px-6 py-12 text-center" style={{ color: 'rgb(var(--text-secondary))' }}>
            No boards available yet.
          </div>
        ) : (
          <div className="divide-y" style={{ borderColor: 'rgb(var(--border))' }}>
            {boards?.map((board) => (
              <Link
                key={board.id}
                to="/board/$slug"
                params={{ slug: board.slug }}
                className="board-link block px-6 py-4 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold" style={{ color: 'rgb(var(--text-primary))' }}>
                      /{board.slug}/
                    </h3>
                    <p className="mt-1" style={{ color: 'rgb(var(--text-secondary))' }}>
                      {board.description}
                    </p>
                  </div>
                  <div className="flex items-center space-x-4">
                    <div className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>
                      {board.topic_count} topic{board.topic_count !== 1 ? 's' : ''}
                    </div>
                    <div style={{ color: 'rgb(var(--text-secondary))' }}>
                      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                      </svg>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export const Route = createFileRoute('/')({
  component: Boards,
})
