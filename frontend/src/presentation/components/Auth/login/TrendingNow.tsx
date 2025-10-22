
import React, { useState } from "react"
import { IoMdClose } from "react-icons/io";
import { motion } from "framer-motion";
import { movieList } from "../../../constants/LoginMovies";

interface FilmProps {
    id?:number
    link:string
    desc:string  
    categories : 
    {
        title : string
        year : number
        age_rating : string
        genre : string
    }
    setDisplayInfo:React.Dispatch<React.SetStateAction<boolean>>
}


const FilmInfoModal = ({link , desc , categories , setDisplayInfo}:FilmProps)=>{

    const handleDesactiveModal  =()=>{
        setDisplayInfo(false);
    }

    return (
    <div className="fixed top-0 left-0 h-screen w-screen bg-gradient-to-br from-black/70 via-zinc-900/80 to-gray-800/70 backdrop-blur-sm z-[999999] flex justify-center items-center">
        <div className={`relative text-white`}>
                <div className="absolute right-0 top-0">
                    <button onClick={handleDesactiveModal}>
                        <IoMdClose/>
                    </button>
                </div>
                <section>
                    <img src={link} />
                    <div className="flex gap-2">
                        <span>{categories?.title}</span>
                        <span>{categories?.year}</span>
                        <span>{categories?.age_rating}</span>
                        <span>{categories?.genre}</span>
                    </div>
                    <p>{desc}</p>
                </section>
            </div>
        </div>
    )
}


const Film = ({id , link , desc , categories}:FilmProps)=>{

    const [hover,setHover]=useState(false);
    const [displayInfo,setDisplayInfo]=useState(false);
    const handleActiveModal = ()=>{
        setDisplayInfo(true);
    }

    return (
        <>
            <div className="relative hover:scale-125 hover:cursor-pointer active:scale-125" onClick={handleActiveModal}>
                <span className="absolute  text-4xl text-white text-outline ">{id}</span>
                <img src={link} className="h-full max-w-[200px]"/>
            </div>
            {displayInfo && <FilmInfoModal link={link} desc={desc} categories={categories} setDisplayInfo={setDisplayInfo}/>}
        </>
    )
}




export default function TrendingNow() {
  const [index, setIndex] = useState(0); // البداية من الفيلم الأول
  const step = 1; // عدد الأفلام اللي تتحرك فمرة وحدة

  const handlePrev = () => {
    setIndex((prev) => Math.max(prev - step, 0));
  };

  const handleNext = () => {
    setIndex((prev) => Math.min(prev + step, movieList.length - 1) );
  };

  return (
    <div className="relative w-full overflow-hidden">
      <h1 className="text-3xl font-bold mb-4">Trending Now</h1>
      <div className="flex items-center">
        <button
          onClick={handlePrev}
          className="z-10 bg-gray-800 text-white px-3 py-2 rounded hover:bg-gray-700"
        >
          {"<"}
        </button>
        <div className="overflow-hidden flex-1 relative">
          <motion.div
            className="flex gap-4"
            animate={{ x: -index * 260 }} // 260px = width + gap تقريباً لكل فيلم
            transition={{ type: "spring", stiffness: 200, damping: 30 }}
          >
            {movieList.map((film) => (
              <Film
                    key={film.id}
                    id={film.id}
                    link={film.imageLink}
                    desc={film.desc}
                    categories={film.categories} setDisplayInfo={function (value: React.SetStateAction<boolean>): void {
                        throw new Error("Function not implemented.");
                    } }              />
            ))}
          </motion.div>
        </div>
        <button
          onClick={handleNext}
          className="z-10 bg-gray-800 text-white px-3 py-2 rounded hover:bg-gray-700"
        >
          {">"}
        </button>
      </div>
    </div>
  );
}

