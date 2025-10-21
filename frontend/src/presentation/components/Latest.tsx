import React, { useState } from 'react';
// import { FaStar } from 'react-icons/fa'; // Uncomment if you want to add star ratings

const FilmCard = ({ film }) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className="relative w-72 h-96 rounded-lg overflow-hidden shadow-lg transform transition-transform duration-300 hover:scale-105 cursor-pointer"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Film Poster */}
      <img
        src={film.posterUrl}
        alt={film.title}
        className="w-full h-full object-cover"
      />

      {/* Tags (Types) - Always Visible */}
      <div className="absolute top-2 left-2 flex flex-wrap gap-1">
        {film.tags.map((tag, index) => (
          <span
            key={index}
            className="bg-red-600 text-white text-xs font-semibold px-2 py-1 rounded-full opacity-90"
          >
            {tag}
          </span>
        ))}
      </div>

      {/* Overlay Info on Hover */}
      <div
        className={`absolute inset-0 bg-gradient-to-t from-black via-black/70 to-transparent flex flex-col justify-end p-4 text-white transition-opacity duration-300 ${
          isHovered ? 'opacity-100' : 'opacity-0'
        }`}
      >
        <h3 className="text-xl font-bold mb-1">{film.title}</h3>
        <p className="text-sm mb-2 line-clamp-3"> {/* line-clamp for description */}
          {film.description}
        </p>
        <div className="flex items-center text-sm">
          {/* Uncomment if you want to add star ratings */}
          {/* <FaStar className="text-yellow-400 mr-1" /> */}
          <span>{film.year} | {film.runtime}</span>
          {/* You could add more info here, e.g., | {film.rating}/5 */}
        </div>
      </div>
    </div>
  );
};

const Latest = () => {
  // Mock data for your films
  const films = [
    {
      id: 1,
      title: 'The Enigma Code',
      posterUrl: '/films/prisonBreak.jpeg', // Replace with a relevant image URL
      description: 'A brilliant detective must decipher a series of cryptic clues to prevent a global catastrophe. Full of twists and turns, this thriller will keep you on the edge of your seat.',
      tags: ['Thriller', 'Mystery', 'Action'],
      year: 2023,
      runtime: '1h 55m',
      // rating: 4.5,
    },
    {
      id: 2,
      title: 'Whispers of the Galaxy',
      posterUrl: '/films/prisonBreak.jpeg', // Replace with a relevant image URL
      description: 'In a distant future, a lone explorer uncovers an ancient secret that could change the fate of the universe. An epic space opera with breathtaking visuals and emotional depth.',
      tags: ['Sci-Fi', 'Adventure', 'Fantasy'],
      year: 2022,
      runtime: '2h 10m',
      // rating: 4.8,
    },
    // You can add more films here if needed
  ];

  return (
    <div className="w-full bg-black py-16 px-4 md:px-8 lg:px-16 border-b-8 border-[#222]">
      <h2 className="text-white text-3xl md:text-4xl font-bold mb-10 text-center">
        Latest & Trending (Fakeflix Originals)
      </h2>
      <div className="flex justify-center items-center flex-wrap gap-8">
        {films.map((film) => (
          <FilmCard key={film.id} film={film} />
        ))}
      </div>
    </div>
  );
};

export default Latest;