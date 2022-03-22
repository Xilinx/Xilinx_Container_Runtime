AC_INIT([xilinx-container-runtime],[1.0.0],[yahuil@xilinx.com])

AC_CHECK_PROG(GO_CHECK,go,yes)
if test x"$GO_CHECK" != x"yes" ; then
    AC_MSG_ERROR([Please install golang to build the binary.])
fi


AC_CHECK_PROG(RUNC_CHECK,runc,yes)
if test x"$RUNC_CHECK" != x"yes" ; then
    AC_MSG_WARN([RunC is required to start a container.])
fi