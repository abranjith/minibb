import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/debug')({
  component: () => {
    console.log('🐛 Debug page loaded successfully!')
    
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold mb-4">Debug Page</h1>
        <p>If you can see this, routing is working!</p>
        <p>Current URL: {window.location.href}</p>
        <p>Time: {new Date().toLocaleString()}</p>
        <div>
          <Link to="/board/general">Go to /board/general</Link>
          <br />
          <Link to="/board/general/topic/1">Go to /board/general/topic/1</Link>
        </div>
      </div>
    )
  },
})
