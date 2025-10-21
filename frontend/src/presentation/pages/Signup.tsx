import { useState } from 'react'
import { FaPlay , FaPause } from "react-icons/fa";
import TrendingNow from '../components/Auth/login/TrendingNow';
import Register from '../components/Auth/Register';
import './Signup.css'
import { useTranslation } from 'react-i18next';

const Signup = () => {
  const [play,setPlay] = useState<boolean>(false);
  const {t} = useTranslation();
  const handlePlayClick = ()=>{
    setPlay(!play);
  }
  return (
    <div>
        <div className='relative animate-scroll-x h-[40vh] flex flex-col  gap-6 justify-center items-center'>
            <p className='font-bold text-white max-w-[200px] md:max-w-[600px] text-center text-3xl md:text-6xl'>
                {t('title')}
            </p>
            <Register/>
            <button className='absolute bottom-4 right-4 bg-black text-white p-2 rounded-full'
            onClick={handlePlayClick}>
                {!play ? <FaPause/> : <FaPlay/> }
            </button>
        </div>
        <TrendingNow/>
    </div>
  )
}

export default Signup