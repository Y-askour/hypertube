import { movieList } from "@/presentation/constants/LoginMovies"
import React, { useState } from "react"
import { IoMdClose } from "react-icons/io";

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
                <span className="absolute  text-3xl ">{id}</span>
                <img src={link}/>
            </div>
            {displayInfo && <FilmInfoModal link={link} desc={desc} categories={categories} setDisplayInfo={setDisplayInfo}/>}
        </>
    )
}


const TrendingNow = () => {
  
  
  return (
    <div>
        <h1>Trending Now</h1>
        <div className="flex gap-4">
            <button></button>
            {movieList.map((film)=>
                <Film id={film?.id} link={film.imageLink} desc={film.desc} categories={film.categories}/>
            )}
            <button></button>
        </div>
    </div>
  )
}

export default TrendingNow