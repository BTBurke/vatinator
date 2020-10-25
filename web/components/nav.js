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
