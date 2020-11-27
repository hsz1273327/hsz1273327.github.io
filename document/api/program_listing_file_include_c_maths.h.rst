
.. _program_listing_file_include_c_maths.h:

Program Listing for File c_maths.h
==================================

|exhale_lsh| :ref:`Return to documentation for file <file_include_c_maths.h>` (``include\c_maths.h``)

.. |exhale_lsh| unicode:: U+021B0 .. UPWARDS ARROW WITH TIP LEFTWARDS

.. code-block:: cpp

   /***************************************************************************************
    * This file is dedicated to the public domain.  If your jurisdiction requires a       *
    * specific license:                                                                   *
    *                                                                                     *
    * Copyright (c) Stephen McDowell, 2017-2019                                           *
    * License:      CC0 1.0 Universal                                                     *
    * License Text: https://creativecommons.org/publicdomain/zero/1.0/legalcode           *
    **************************************************************************************/
   #ifndef C_MATHS_H
   #define C_MATHS_H
   
   #if defined(__cplusplus)
       extern "C" {
   #endif
   
   int cm_add(int a, int b);
   
   int cm_sub(int a, int b);
   
   #if defined(__cplusplus)
       } // extern "C"
   #endif
   
   #endif // C_MATHS_H
