import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import SignIn from '../pages/Login';


const languages = ['English' , 'French'];

function NavBar() {
  const [selectedLanguage,setSelectedLanguage] = useState();
  const [open,setOpen]=useState<boolean>(false);
  const { i18n , t } = useTranslation();
  const changeLanguage = (lng:string) => {
    i18n.changeLanguage(lng);
  }

  const handleSelectLanguage =(e:any)=>{
    const {value} = e.target;
    setSelectedLanguage(value);
  }

  useEffect(()=>{
    const lng = selectedLanguage === 'English' ? 'en' : selectedLanguage === 'French' ? 'fr' : 'en'
    localStorage.setItem('i18nextLng',lng);
    changeLanguage(lng)
  },[selectedLanguage])

  return (
    <div className='flex justify-between '>
        <h1 className='text-red-500 font-bold text-2xl'>Hypertube</h1>
        <div className='flex gap-2'>
          <button onClick={()=>setOpen(true)} className='bg-green-300 py-1 px-4 font-bold text-sm rounded-full bg-white hover:bg-gray-300'>{t('Sign in')}</button>
          <select value={selectedLanguage} onChange={handleSelectLanguage} className='px-2 py-1 rounded-md bg-black text-white border-[1px] border-white'>
            {languages.map((lang)=><option key={lang}>{lang}</option>)}
          </select>
        {
          open && <SignIn setOpen={setOpen}/>
        }
        </div>
    </div>
  )
}

export default NavBar