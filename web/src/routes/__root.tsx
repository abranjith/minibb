import { createRootRoute, Outlet, Link } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'
import { ThemeToggle } from '../components/ThemeToggle'

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="min-h-screen themed-bg">
        <header className="themed-header shadow-md">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-center relative">
            <Link 
              to="/"
              className="text-2xl font-bold hover:opacity-80 transition-opacity cursor-pointer inline-block no-underline header-text"
              activeProps={{
                className: "text-2xl font-bold opacity-80 cursor-pointer inline-block no-underline header-text"
              }}
            >
              MiniBB
            </Link>
            <div className="absolute right-1">
              <ThemeToggle />
            </div>
          </div>
        </header>
        <main className="max-w-7xl mx-auto px-4 py-6 themed-text">
          <Outlet />
        </main>
      </div>
      <TanStackRouterDevtools />
    </>
  ),
})
