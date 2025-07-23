import { createFileRoute, Outlet } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { hasUnreadPosts } from '../../utils/readStatus'

interface IBoard {
  id: number
  slug: string
  description: string
}

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

async function fetchBoard(slug: string): Promise<IBoard> {
  const response = await fetch(`/api/board/${slug}`)
  if (!response.ok) {
    throw new Error('Failed to fetch board')
  }
  return response.json()
}

async function fetchTopics(boardId: number): Promise<ITopic[]> {
  const response = await fetch(`/api/topics/board/${boardId}`)
  if (!response.ok) {
    throw new Error('Failed to fetch topics')
  }
  const data = await response.json()
  return data.topics
}

function Board() {
  const { slug } = Route.useParams()

  const { data: board, isLoading: boardLoading } = useQuery({
    queryKey: ['board', slug],
    queryFn: () => fetchBoard(slug),
  })

  const { data: topics, isLoading: topicsLoading } = useQuery({
    queryKey: ['topics', board?.id],
    queryFn: () => fetchTopics(board!.id),
    enabled: !!board,
  })

  if (boardLoading) {
    return <div className="text-center py-8">Loading board...</div>
  }

  return (
    <div className="space-y-6">
      {/* Back button */}
      <div className="flex items-center gap-2">
        <Link
          to="/"
          className="inline-flex items-center gap-2 text-blue-600 hover:text-blue-800 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to Boards
        </Link>
      </div>

      <Outlet/>

      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold text-primary">/{board?.slug}/</h2>
          <p className="text-secondary mt-2">{board?.description}</p>
        </div>
        <Link
          to="/board/$slug/new-topic"
          params={{ slug }}
          className="btn-primary"
        >
          New Topic
        </Link>
      </div>

      {topicsLoading ? (
          <div className="text-center py-8">Loading topics...</div>
      ) : (
        <div className="themed-bg-secondary rounded-lg shadow-md overflow-hidden">
          <div className="px-6 py-4 themed-bg border-b themed-border">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 text-sm font-medium text-secondary">
              <div className="md:col-span-2">Topic</div>
              <div className="hidden md:block">Author</div>
              <div className="hidden md:block">Replies</div>
            </div>
          </div>
          
          {topics?.length === 0 ? (
            <div className="px-6 py-12 text-center text-secondary">
              No topics yet. Be the first to start a discussion!
            </div>
          ) : (
            <div className="divide-y themed-border">
              {topics?.map((topic) => (
                <div key={topic.id}>
                  <Link
                    key={`link-${topic.id}`}
                    to="/board/$slug/topic/$topicId"
                    params={{ slug, topicId: topic.id.toString() }}
                    className="block topic-hover transition-colors no-underline"
                    style={{ color: 'inherit' }}
                  >
                    <div className="px-6 py-4">
                      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div className="md:col-span-2 flex items-center gap-2">
                          <h3 className="font-medium text-primary hover:text-accent">
                            {topic.title}
                          </h3>
                          {hasUnreadPosts(topic.id, topic.last_post_id) && (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium new-badge">
                              New
                            </span>
                          )}
                        </div>
                        <div className="text-sm text-secondary">
                          {topic.author}
                        </div>
                        <div className="text-sm text-secondary flex items-center justify-between">
                          <span>{topic.post_count} posts</span>
                          <span className="text-xs text-secondary">
                            {new Date(topic.pub_date).toLocaleDateString()}
                          </span>
                        </div>
                      </div>
                    </div>
                  </Link>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export const Route = createFileRoute('/board/$slug')({
  component: Board,
})
