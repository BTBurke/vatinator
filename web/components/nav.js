import Link from 'next/link'

// const links = [
//   { href: 'https://github.com/vercel/next.js', label: 'Batches' },
// ]

const links = []

export default function Nav() {
  return (
    <nav>
      <ul className="flex justify-around lg:justify-between items-center py-8 px-4">
        <li>
          <img style={{height: "70px"}} src="/vatinatorAsPath.svg" />
        </li>
        {/* <ul className="flex justify-between items-center space-x-4">
          {links.map(({ href, label }) => (
            <li key={`${href}${label}`}>
              <a href={href} className="text-white no-underline">
                {label}
              </a>
            </li>
          ))}
        </ul> */}
      </ul>
    </nav>
  )
}

export function NavHome() {
  return (
  <nav>
      <ul className="flex justify-around lg:justify-between items-center py-8 px-4">
        <li>
          <img style={{height: "70px"}} src="/vatinatorAsPath.svg" />
        </li>
        <li className="hidden md:block">
          <a className="mr-4 bg-primary h-12 text-white px-full py-2 px-6 rounded-md font-bold border border-white" href="/login">Login</a>
          <a className="mr-4 bg-primary h-12 text-white px-full py-2 px-6 rounded-md font-bold border border-white" href="/create">Create account</a>
        </li>
      </ul>
    </nav>
  )
}
