import { ReactNode } from "react";
import Footer from "./Footer";
import Navbar from "./Navbar";
import styled from "styled-components";

export default function Page({children}: { children: ReactNode}) {
  return (
    <Wrapper>
      <Navbar />
      <Middle>
        {children}
      </Middle>
      <Footer />
    </Wrapper>
  )
}

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
`;

const Middle = styled.div`
  flex: 1;
  background-color: pink;
  color: black;
`;