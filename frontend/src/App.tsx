import FAQSection from "./presentation/components/FAQSection";
import Latest from "./presentation/components/Latest";
import Footer from "./presentation/layout/Footer";
import NavBar from "./presentation/layout/NavBar";
import Signup from "./presentation/pages/Signup";


function App() {

  return (
    <div className="min-h-screen bg-gradient-to-r from-[#0e090c] to-[#e50914] pt-4 flex flex-col gap-4">
      <NavBar/>
      <Signup/>
      <Latest/>
      <FAQSection/>
      <Footer/>
    </div>
    
  );
}

export default App;
