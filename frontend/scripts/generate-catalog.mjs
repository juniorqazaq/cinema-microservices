/**
 * Writes src/data/{trending,upcoming,nowPlaying,topRated}.json
 * Poster/backdrop paths: TMDB file paths only (see scripts/catalog-tmdb-paths.mjs).
 * Run: node scripts/generate-catalog.mjs (from frontend/)
 */
import { writeFileSync, mkdirSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { TMDB_PATHS } from './catalog-tmdb-paths.mjs'

const __dirname = dirname(fileURLToPath(import.meta.url))
const dataDir = join(__dirname, '..', 'src', 'data')

const castPool = [
  ['Timothée Chalamet', 'Paul Atreides'],
  ['Zendaya', 'Chani'],
  ['Rebecca Ferguson', 'Lady Jessica'],
  ['Oscar Isaac', 'Duke Leto'],
  ['Florence Pugh', 'Princess Irulan'],
  ['Austin Butler', 'Feyd-Rautha'],
  ['Cillian Murphy', 'Robert Oppenheimer'],
  ['Emily Blunt', 'Kitty Oppenheimer'],
  ['Matt Damon', 'Leslie Groves'],
  ['Robert Downey Jr.', 'Lewis Strauss'],
  ['Tom Cruise', 'Ethan Hunt'],
  ['Hayley Atwell', 'Grace'],
  ['Ana de Armas', 'Dani Miranda'],
  ['Ryan Gosling', 'Colt Seavers'],
  ['Emily Blunt', 'Jody Moreno'],
  ['Pedro Pascal', 'Tom Ryder'],
  ['Paul Giamatti', 'Paul Hunham'],
  ['Da’Vine Joy Randolph', 'Mary Lamb'],
  ['Dominic Sessa', 'Angus Tully'],
  ['Lily Gladstone', 'Mollie Burkhart'],
  ['Leonardo DiCaprio', 'Ernest Burkhart'],
  ['Robert De Niro', 'William Hale'],
  ['Jesse Plemons', 'Tom White'],
  ['Margot Robbie', 'Barbie'],
  ['Ryan Gosling', 'Ken'],
  ['America Ferrera', 'Gloria'],
  ['Michael Cera', 'Allan'],
  ['Ke Huy Quan', 'Ode to Joy narrator'],
]

function reviewsFor(title, baseRating) {
  return [
    {
      author: 'Rotten Tomatoes–style critic',
      rating: Math.min(10, baseRating + 0.3),
      content: `${title} delivers the spectacle audiences expect while keeping character stakes grounded. A rare blockbuster that earns its runtime.`,
    },
    {
      author: 'IMDb-style viewer',
      rating: Math.min(10, baseRating - 0.2),
      content: `Saw it twice. Sound design and pacing are excellent; a few plot beats feel rushed but the cast chemistry carries the third act.`,
    },
  ]
}

function buildCast(seed) {
  const out = []
  for (let i = 0; i < 4; i++) {
    const idx = (seed + i * 7) % castPool.length
    const [name, character] = castPool[idx]
    out.push({
      name,
      character,
      profile_path: `https://i.pravatar.cc/300?u=${encodeURIComponent(name + seed)}`,
    })
  }
  return out
}

function movie(row) {
  const {
    id,
    title,
    overview,
    vote_average,
    vote_count,
    popularity,
    release_date,
    genre_ids,
    runtime,
    status,
    tagline,
    trailer,
    poster_path,
    backdrop_path,
    adult = false,
    original_language = 'en',
  } = row

  return {
    id,
    title,
    original_title: title,
    overview,
    poster_path,
    backdrop_path,
    release_date,
    vote_average,
    vote_count,
    popularity,
    genre_ids,
    adult,
    original_language,
    runtime,
    status,
    tagline,
    reviews: reviewsFor(title, vote_average),
    cast: buildCast(id),
    trailer,
  }
}

const trailers = [
  'Way9Dexny3w',
  '6RRLLCv2srQ',
  'eWc47_ocZb0',
  'd9MyW72ELq0',
  'L9W88mbtZxw',
  'm8GE3m5CEN8',
  'pBk4NYhWNMM',
  'uYPbbksJxIg',
  'h74AXqw4Opc',
  'otNh9gTtfXU',
  'pVsOZXrO2o0',
  'wKiI0RvKgSe',
  'Shah8AN253Y',
  'uqmAj8dox64',
  'yAZxx8t9zig',
]

function yt(i) {
  return `https://www.youtube.com/watch?v=${trailers[i % trailers.length]}`
}

const trendingRows = [
  [101, 'Dune: Part Two', 'Paul unites with the Fremen on a warpath of revenge against those who destroyed his family.', 8.6, 420000, 9800, '2024-03-01', [878, 12], 166, 'Released', 'Long live the fighters.'],
  [102, 'Oppenheimer', 'The story of American scientist J. Robert Oppenheimer and his role in developing the atomic bomb.', 8.4, 610000, 9200, '2023-07-21', [36, 18], 180, 'Released', 'The world forever changes.'],
  [103, 'Mission: Impossible – Dead Reckoning', 'Ethan Hunt and the IMF team must track down a terrifying new weapon.', 7.7, 310000, 7800, '2023-07-12', [28, 12], 163, 'Released', 'We all share the same fate.'],
  [104, 'Barbie', 'Barbie suffers a crisis that leads her to question her world and her existence.', 6.9, 890000, 8600, '2023-07-21', [35, 14], 114, 'Released', 'She’s everything. He’s just Ken.'],
  [105, 'Killers of the Flower Moon', 'Members of the Osage tribe are murdered under mysterious circumstances.', 7.7, 195000, 5400, '2023-10-20', [36, 80], 206, 'Released', 'Greed is an animal.'],
  [106, 'The Holdovers', 'A cranky history teacher at a remote prep school is forced to remain on campus over the holidays.', 7.9, 142000, 4100, '2023-10-27', [35, 18], 133, 'Released', 'Discomfort & joy.'],
  [107, 'Poor Things', 'Brought back to life by an eccentric scientist, a young woman runs off on a whirlwind adventure.', 7.8, 178000, 5200, '2023-12-08', [878, 35], 141, 'Released', 'She’s like nothing you’ve ever seen.'],
  [108, 'John Wick: Chapter 4', 'John Wick uncovers a path to defeating The High Table.', 7.7, 290000, 7100, '2023-03-24', [28, 53], 169, 'Released', 'No way back. One way out.'],
  [109, 'Spider-Man: Across the Spider-Verse', 'Miles Morales catapults across the Multiverse to stop a threat to all realities.', 8.4, 520000, 8800, '2023-06-02', [16, 12], 140, 'Released', 'Be greater. Together.'],
  [110, 'Guardians of the Galaxy Vol. 3', 'Still reeling from the loss of Gamora, Peter Quill must rally his team.', 7.9, 380000, 6700, '2023-05-05', [12, 878], 150, 'Released', 'It’s time to face the music.'],
  [111, 'Past Lives', 'Childhood sweethearts are reunited for one fateful week as they contemplate what could have been.', 8.0, 112000, 3900, '2023-06-02', [10749, 18], 105, 'Released', 'What if…'],
  [112, 'The Batman', 'Batman ventures into Gotham’s underworld when a sadistic killer leaves behind a trail of cryptic clues.', 7.8, 610000, 7900, '2022-03-04', [80, 28], 176, 'Released', 'Unmask the truth.'],
  [113, 'Everything Everywhere All at Once', 'A middle-aged Chinese immigrant is swept into an insane adventure.', 7.8, 480000, 8200, '2022-03-25', [28, 878], 139, 'Released', 'The universe is vast.'],
  [114, 'Top Gun: Maverick', 'After thirty years, Maverick is still pushing the envelope as a top naval aviator.', 8.2, 720000, 9100, '2022-05-27', [28, 18], 131, 'Released', 'Feel the need for speed.'],
  [115, 'Avatar: The Way of Water', 'Jake Sully lives with his newfound family on the extrasolar moon Pandora.', 7.5, 620000, 8500, '2022-12-16', [878, 12], 192, 'Released', 'Return to Pandora.'],
]

const upcomingRows = [
  [201, 'Thunderbolts*', 'A team of antiheroes embark on missions for the government.', 0, 12000, 7200, '2026-07-17', [28, 12], 0, 'Post-production', 'Not heroes. Not villains.'],
  [202, 'The Fantastic Four: First Steps', 'Marvel’s first family faces their greatest challenge yet.', 0, 8900, 6800, '2026-07-25', [878, 12], 0, 'Post-production', 'Prepare 4 launch.'],
  [203, 'Superman', 'Superman balances his Kryptonian heritage with his human upbringing.', 0, 15000, 9100, '2026-07-11', [28, 14], 0, 'Post-production', 'Look up.'],
  [204, 'Mission: Impossible – The Final Reckoning', 'Ethan Hunt faces his deadliest mission yet.', 0, 22000, 8400, '2026-05-22', [28, 12], 0, 'Post-production', 'Nothing ends well.'],
  [205, 'Jurassic World: Rebirth', 'A new era of survival begins.', 0, 11000, 6200, '2026-07-02', [12, 28], 0, 'Post-production', 'Life finds a way.'],
  [206, 'F1', 'A Formula 1 driver returns from retirement.', 0, 9500, 5800, '2026-06-27', [18, 28], 0, 'Post-production', 'Racing is life.'],
  [207, 'Weapons', 'A horror mystery surrounding small-town disappearances.', 0, 4200, 3100, '2026-08-08', [27, 9648], 0, 'Post-production', 'Everyone has a weapon.'],
  [208, 'The Bride!', 'A gothic romance reimagined for modern audiences.', 0, 3100, 2800, '2026-03-06', [14, 10749], 0, 'Post-production', 'Love never dies.'],
  [209, 'Coyote vs. Acme', 'Wile E. Coyote sues the Acme Corporation.', 0, 28000, 7600, '2026-06-12', [16, 35], 0, 'Post-production', 'Order in the court.'],
  [210, 'Anaconda', 'A team ventures deep into the rainforest.', 0, 5600, 3400, '2026-12-25', [28, 27], 0, 'Pre-production', 'Fear coils.'],
  [211, 'Project Hail Mary', 'A lone astronaut races to save humanity.', 0, 18000, 8900, '2026-03-20', [878, 12], 0, 'Filming', 'Science is the only hope.'],
  [212, 'Mortal Kombat II', 'Earthrealm’s champions fight for survival.', 0, 14000, 5500, '2026-10-24', [28, 14], 0, 'Post-production', 'Finish them.'],
  [213, 'Now You See Me: Now You Don’t', 'The Four Horsemen return for their most dangerous illusion yet.', 0, 7200, 4100, '2026-11-13', [80, 9648], 0, 'Pre-production', 'Look closer.'],
  [214, 'Hoppers', 'An animated adventure across impossible worlds.', 0, 2100, 2200, '2026-03-06', [16, 12], 0, 'Post-production', 'Hop into the unknown.'],
  [215, 'The Sheep Detectives', 'Two woolly investigators crack barnyard mysteries.', 0, 900, 1500, '2026-09-18', [16, 35], 0, 'In development', 'Ewe won’t believe it.'],
]

const nowPlayingRows = [
  [301, 'Sinners', 'Twin brothers return to their hometown to escape the past—only to face something far older.', 7.8, 89000, 6200, '2025-04-18', [27, 36], 137, 'Released', 'Pray for dawn.'],
  [302, 'A Minecraft Movie', 'Four misfits are pulled through a portal into the Overworld.', 6.2, 210000, 7800, '2025-04-04', [12, 35], 101, 'Released', 'Be anything.'],
  [303, 'The Accountant 2', 'Christian Wolff applies his skills to a deadly new case.', 6.9, 45000, 4800, '2025-04-25', [28, 80], 132, 'Released', 'Numbers don’t lie.'],
  [304, 'Thunderbolts*', 'Antiheroes assemble for a covert mission.', 7.1, 67000, 7100, '2025-05-02', [28, 12], 126, 'Released', 'Controlled chaos.'],
  [305, 'Fight or Flight', 'A passenger must protect a mysterious package at 30,000 feet.', 6.4, 12000, 3600, '2025-05-09', [28, 53], 98, 'Released', 'No parachute.'],
  [306, 'Final Destination: Bloodlines', 'Death’s design returns with a vengeance.', 6.8, 34000, 5200, '2025-05-16', [27, 53], 109, 'Released', 'You can’t cheat fate.'],
  [307, 'Lilo & Stitch', 'A lonely girl adopts a chaotic alien “dog.”', 7.0, 78000, 6400, '2025-05-23', [10749, 16], 108, 'Released', 'Ohana means family.'],
  [308, 'Karate Kid: Legends', 'Martial arts legacies collide across generations.', 6.9, 52000, 5100, '2025-05-30', [28, 18], 124, 'Released', 'Two disciplines. One destiny.'],
  [309, 'Ballerina', 'A young assassin seeks revenge in the John Wick universe.', 7.2, 41000, 5800, '2025-06-06', [28, 53], 125, 'Released', 'Dance with death.'],
  [310, 'How to Train Your Dragon', 'A live-action retelling of the bond between Viking and dragon.', 7.4, 96000, 6900, '2025-06-13', [12, 14], 126, 'Released', 'Fly together.'],
  [311, 'Elio', 'A boy is beamed up to represent Earth to alien diplomats.', 6.6, 28000, 4400, '2025-06-20', [16, 878], 98, 'Released', 'Wrong planet. Right kid.'],
  [312, 'F1', 'High stakes on and off the track.', 7.6, 55000, 7200, '2025-06-27', [18, 28], 155, 'Released', 'Full throttle.'],
  [313, 'Superman', 'Hope returns when Krypton’s last son chooses Earth.', 7.5, 125000, 8800, '2025-07-11', [28, 14, 878], 129, 'Released', 'Truth. Justice. Tomorrow.'],
  [314, 'Jurassic World: Rebirth', 'New island. New predators. Same survival instinct.', 6.9, 88000, 7400, '2025-07-02', [12, 28], 134, 'Released', 'Evolution bites.'],
  [315, 'The Fantastic Four: First Steps', 'Marvel’s first family steps into the unknown.', 7.3, 102000, 8100, '2025-07-25', [878, 12], 115, 'Released', '4 ever.'],
]

const topRatedRows = [
  [401, 'The Shawshank Redemption', 'Two imprisoned men bond over years, finding solace and eventual redemption.', 9.3, 3100000, 12000, '1994-09-23', [18, 80], 142, 'Released', 'Fear can hold you prisoner. Hope can set you free.'],
  [402, 'The Godfather', 'The aging patriarch of an organized crime dynasty transfers control to his reluctant son.', 9.2, 2100000, 11000, '1972-03-24', [80, 18], 175, 'Released', 'An offer you can’t refuse.'],
  [403, 'The Dark Knight', 'Batman faces the Joker, who throws Gotham into anarchy.', 9.0, 3200000, 10500, '2008-07-18', [28, 80, 53], 152, 'Released', 'Why so serious?'],
  [404, 'The Godfather Part II', 'The early life and career of Vito Corleone in 1920s New York.', 9.0, 1300000, 9800, '1974-12-20', [80, 18], 202, 'Released', 'All the power on Earth.'],
  [405, '12 Angry Men', 'A jury holdout attempts to prevent a miscarriage of justice.', 9.0, 820000, 7200, '1957-04-10', [18, 80], 96, 'Released', 'Life is in their hands.'],
  [406, 'Schindler’s List', 'In German-occupied Poland, Oskar Schindler gradually becomes concerned for his Jewish workers.', 9.0, 1400000, 8600, '1993-12-15', [36, 18, 10752], 195, 'Released', 'Whoever saves one life, saves the world entire.'],
  [407, 'The Lord of the Rings: The Return of the King', 'Gandalf and Aragorn lead the World of Men against Sauron’s army.', 8.9, 1900000, 9900, '2003-12-17', [12, 14, 28], 201, 'Released', 'The eye of the enemy is moving.'],
  [408, 'Pulp Fiction', 'The lives of two mob hitmen, a boxer, and others intertwine in four tales.', 8.9, 2100000, 9400, '1994-10-14', [80, 53], 154, 'Released', 'Guns. Dialogue. Royale with cheese.'],
  [409, 'The Good, the Bad and the Ugly', 'A bounty hunting scam joins two men in an uneasy alliance.', 8.8, 760000, 7800, '1966-12-23', [37, 28], 178, 'Released', 'For three men the Civil War wasn’t hell.'],
  [410, 'Fight Club', 'An insomniac office worker and a soapmaker form an underground fight club.', 8.8, 2100000, 9000, '1999-10-15', [18, 53], 139, 'Released', 'Mischief. Mayhem. Soap.'],
  [411, 'Forrest Gump', 'The presidencies of Kennedy and Johnson, Vietnam, and more through an Alabama man’s eyes.', 8.8, 2100000, 8800, '1994-07-06', [18, 10749, 10752], 142, 'Released', 'The world will never be the same.'],
  [412, 'Inception', 'A thief who steals corporate secrets through dream-sharing is offered a chance at redemption.', 8.8, 2400000, 10200, '2010-07-16', [28, 878, 9648], 148, 'Released', 'Your mind is the scene of the crime.'],
  [413, 'The Lord of the Rings: The Fellowship of the Ring', 'A meek Hobbit and eight companions set out to destroy the One Ring.', 8.8, 2000000, 9700, '2001-12-19', [12, 14, 28], 178, 'Released', 'One ring to rule them all.'],
  [414, 'Star Wars: Episode V – The Empire Strikes Back', 'After the Rebels are brutally overpowered, Luke Skywalker begins Jedi training.', 8.7, 1400000, 9200, '1980-05-21', [12, 28, 878], 124, 'Released', 'The adventure continues…'],
  [415, 'The Matrix', 'A computer hacker learns about the true nature of reality.', 8.7, 2000000, 10100, '1999-03-31', [28, 878], 136, 'Released', 'Reality is a choice.'],
]

function rowsToMovies(rows, offset) {
  return rows.map((r, i) => {
    const id = r[0]
    const pair = TMDB_PATHS[id] ?? TMDB_PATHS[String(id)]
    const [poster_path, backdrop_path] = pair
    return movie({
      id: r[0],
      title: r[1],
      overview: r[2],
      vote_average: r[3],
      vote_count: r[4],
      popularity: r[5],
      release_date: r[6],
      genre_ids: r[7],
      runtime: r[8],
      status: r[9],
      tagline: r[10],
      trailer: yt(offset + i),
      poster_path,
      backdrop_path,
    })
  })
}

mkdirSync(dataDir, { recursive: true })

const trending = rowsToMovies(trendingRows, 0)
const upcoming = rowsToMovies(upcomingRows, 20)
const nowPlaying = rowsToMovies(nowPlayingRows, 40)
const topRated = rowsToMovies(topRatedRows, 60)

writeFileSync(join(dataDir, 'trending.json'), JSON.stringify(trending, null, 2))
writeFileSync(join(dataDir, 'upcoming.json'), JSON.stringify(upcoming, null, 2))
writeFileSync(join(dataDir, 'nowPlaying.json'), JSON.stringify(nowPlaying, null, 2))
writeFileSync(join(dataDir, 'topRated.json'), JSON.stringify(topRated, null, 2))

console.log('Wrote 4 catalog files (15 movies each) →', dataDir)
